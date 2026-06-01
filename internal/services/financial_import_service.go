package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
)

// ImportResult summarises one import run.
type ImportResult struct {
	Imported int      `json:"imported"`
	Updated  int      `json:"updated"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// TemplateColumns is the canonical ordered list of CSV/Excel column headers.
// The first row of any uploaded file MUST match these exactly (case-insensitive).
var TemplateColumns = []string{
	"period_end_date",
	"net_profit",
	"total_operating_income",
	"profit_before_tax",
	"total_assets",
	"total_equity",
	"total_liabilities",
	"eps",
	"customer_deposits",
	"loans_and_advances",
}

// TemplateCSV returns the CSV template file content (headers + one example row).
func TemplateCSV() string {
	var b strings.Builder
	b.WriteString(strings.Join(TemplateColumns, ","))
	b.WriteString("\n")
	b.WriteString("2023-12-31,5000000,8000000,6000000,100000000,20000000,80000000,25.50,50000000,30000000\n")
	return b.String()
}

// FinancialImportService handles parsing and upserting financial report rows.
type FinancialImportService struct {
	repo *repositories.FundamentalAnalysisRepository
}

func NewFinancialImportService(repo *repositories.FundamentalAnalysisRepository) *FinancialImportService {
	return &FinancialImportService{repo: repo}
}

// ImportCSV parses a CSV reader and upserts FinancialReport rows for companyProfileID.
func (s *FinancialImportService) ImportCSV(r io.Reader, companyProfileID string) (*ImportResult, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("file must have a header row and at least one data row")
	}
	colIdx, err := mapHeaders(records[0])
	if err != nil {
		return nil, err
	}
	return s.processRows(records[1:], colIdx, companyProfileID)
}

// ImportXLSX parses an Excel file and upserts rows for companyProfileID.
// It reads the first non-empty sheet.
func (s *FinancialImportService) ImportXLSX(r io.Reader, companyProfileID string) (*ImportResult, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("xlsx parse error: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("xlsx file has no sheets")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("xlsx read rows error: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("file must have a header row and at least one data row")
	}
	colIdx, err := mapHeaders(rows[0])
	if err != nil {
		return nil, err
	}
	return s.processRows(rows[1:], colIdx, companyProfileID)
}

func (s *FinancialImportService) processRows(rows [][]string, colIdx map[string]int, companyProfileID string) (*ImportResult, error) {
	result := &ImportResult{}

	for i, row := range rows {
		rowNum := i + 2 // 1-indexed, +1 for header

		// Skip fully empty rows
		allEmpty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			continue
		}

		report, err := rowToReport(row, colIdx, companyProfileID)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", rowNum, err))
			continue
		}

		existing, err := s.repo.GetFinancialReportByCompanyAndPeriod(companyProfileID, report.PeriodEndDate)
		if err == nil && existing != nil {
			// Update
			report.ID = existing.ID
			if err := s.repo.UpdateFinancialReport(report); err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("row %d: update failed: %v", rowNum, err))
				continue
			}
			result.Updated++
		} else {
			// Create
			if err := s.repo.CreateFinancialReport(report); err != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("row %d: insert failed: %v", rowNum, err))
				continue
			}
			result.Imported++
		}
	}

	return result, nil
}

func mapHeaders(header []string) (map[string]int, error) {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}
	if _, ok := idx["period_end_date"]; !ok {
		return nil, fmt.Errorf("missing required column: period_end_date")
	}
	return idx, nil
}

var dateLayouts = []string{"2006-01-02", "02/01/2006", "01/02/2006", "2006/01/02", "02-01-2006"}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised date %q (use YYYY-MM-DD)", s)
}

func parseOptFloat(_ [][]string, colIdx map[string]int, col string, rowSlice []string) *float64 {
	idx, ok := colIdx[col]
	if !ok || idx >= len(rowSlice) {
		return nil
	}
	s := strings.TrimSpace(rowSlice[idx])
	// Strip common currency formatting
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	if s == "" || s == "-" || s == "N/A" || s == "n/a" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

func rowToReport(row []string, colIdx map[string]int, companyProfileID string) (*models.FinancialReport, error) {
	dateIdx, ok := colIdx["period_end_date"]
	if !ok || dateIdx >= len(row) {
		return nil, fmt.Errorf("missing period_end_date value")
	}
	periodEnd, err := parseDate(row[dateIdx])
	if err != nil {
		return nil, err
	}

	// Helper closure to avoid passing redundant args
	optF := func(col string) *float64 {
		return parseOptFloat(nil, colIdx, col, row)
	}

	report := &models.FinancialReport{
		CompanyProfileID:     companyProfileID,
		PeriodEndDate:        periodEnd,
		NetProfit:            optF("net_profit"),
		TotalOperatingIncome: optF("total_operating_income"),
		ProfitBeforeTax:      optF("profit_before_tax"),
		TotalAssets:          optF("total_assets"),
		TotalEquity:          optF("total_equity"),
		TotalLiabilities:     optF("total_liabilities"),
		EPS:                  optF("eps"),
		CustomerDeposits:     optF("customer_deposits"),
		LoansAndAdvances:     optF("loans_and_advances"),
	}
	return report, nil
}
