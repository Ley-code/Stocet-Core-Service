package services

import (
	"errors"
	"strings"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioService struct {
	repo      *repositories.PortfolioRepository
	marketSvc *MarketDataService
}

func NewPortfolioService(repo *repositories.PortfolioRepository, marketSvc *MarketDataService) *PortfolioService {
	return &PortfolioService{repo: repo, marketSvc: marketSvc}
}

type EnrichedHolding struct {
	models.PortfolioHolding
	LatestPrice      float64 `json:"latest_price"`
	MarketValue      float64 `json:"market_value"`
	CostBasis        float64 `json:"cost_basis"`
	UnrealizedPnL    float64 `json:"unrealized_pnl"`
	UnrealizedPnLPct float64 `json:"unrealized_pnl_pct"`
	AllocationPct    float64 `json:"allocation_pct"`
}

type PortfolioSummary struct {
	TotalMarketValue   float64 `json:"total_market_value"`
	TotalCostBasis     float64 `json:"total_cost_basis"`
	TotalUnrealizedPnL float64 `json:"total_unrealized_pnl"`
	TotalPnLPct        float64 `json:"total_unrealized_pnl_pct"`
	HoldingCount       int     `json:"holding_count"`
}

type PortfolioResponse struct {
	Summary  PortfolioSummary  `json:"summary"`
	Holdings []EnrichedHolding `json:"holdings"`
}

type CreateHoldingRequest struct {
	Ticker   string  `json:"ticker"`
	Shares   float64 `json:"shares"`
	BuyPrice float64 `json:"buy_price"`
	BuyDate  string  `json:"buy_date"`
	Notes    string  `json:"notes"`
}

type UpdateHoldingRequest struct {
	Shares   *float64 `json:"shares"`
	BuyPrice *float64 `json:"buy_price"`
	BuyDate  *string  `json:"buy_date"`
	Notes    *string  `json:"notes"`
}

func (s *PortfolioService) GetPortfolio(userID string) (*PortfolioResponse, error) {
	holdings, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	priceMap := s.buildPriceMap()

	enriched := make([]EnrichedHolding, 0, len(holdings))
	var totalMV, totalCost float64

	for _, h := range holdings {
		e := s.enrichOne(h, priceMap)
		enriched = append(enriched, e)
		totalMV += e.MarketValue
		totalCost += e.CostBasis
	}

	for i := range enriched {
		if totalMV > 0 {
			enriched[i].AllocationPct = (enriched[i].MarketValue / totalMV) * 100
		}
	}

	summary := PortfolioSummary{
		TotalMarketValue:   totalMV,
		TotalCostBasis:     totalCost,
		TotalUnrealizedPnL: totalMV - totalCost,
		HoldingCount:       len(enriched),
	}
	if totalCost > 0 {
		summary.TotalPnLPct = ((totalMV - totalCost) / totalCost) * 100
	}

	return &PortfolioResponse{Summary: summary, Holdings: enriched}, nil
}

func (s *PortfolioService) buildPriceMap() map[string]float64 {
	priceMap := make(map[string]float64)
	if s.marketSvc == nil {
		return priceMap
	}
	prices, err := s.marketSvc.GetLatestAll()
	if err != nil {
		return priceMap
	}
	for _, mp := range prices {
		priceMap[mp.TickerSymbol] = mp.Price
	}
	return priceMap
}

func (s *PortfolioService) enrichOne(h models.PortfolioHolding, priceMap map[string]float64) EnrichedHolding {
	cost := h.Shares * h.BuyPrice
	latest := h.BuyPrice

	if price, ok := priceMap[h.Ticker]; ok {
		latest = price
	} else if s.marketSvc != nil && priceMap == nil {
		if mp, err := s.marketSvc.GetLatestByTicker(h.Ticker); err == nil && mp != nil {
			latest = mp.Price
		}
	}

	mv := h.Shares * latest
	pnl := mv - cost
	pnlPct := 0.0
	if cost > 0 {
		pnlPct = (pnl / cost) * 100
	}

	return EnrichedHolding{
		PortfolioHolding: h,
		LatestPrice:      latest,
		MarketValue:      mv,
		CostBasis:        cost,
		UnrealizedPnL:    pnl,
		UnrealizedPnLPct: pnlPct,
	}
}

func (s *PortfolioService) Create(userID string, req CreateHoldingRequest) (*EnrichedHolding, error) {
	ticker := strings.ToUpper(strings.TrimSpace(req.Ticker))
	if ticker == "" {
		return nil, errors.New("ticker is required")
	}
	if req.Shares <= 0 {
		return nil, errors.New("shares must be positive")
	}
	if req.BuyPrice <= 0 {
		return nil, errors.New("buy_price must be positive")
	}

	exists, err := s.repo.ExistsTicker(ticker, userID, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("you already hold this ticker; update the existing position instead")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	h := &models.PortfolioHolding{
		ID:       uuid.New().String(),
		UserID:   uid,
		Ticker:   ticker,
		Shares:   req.Shares,
		BuyPrice: req.BuyPrice,
		Notes:    strings.TrimSpace(req.Notes),
	}
	if bd, err := parseBuyDate(req.BuyDate); err != nil {
		return nil, err
	} else if bd != nil {
		h.BuyDate = bd
	}

	if err := s.repo.Create(h); err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("you already hold this ticker")
		}
		return nil, err
	}

	e := s.enrichOne(*h, nil)
	return &e, nil
}

func (s *PortfolioService) Update(userID, id string, req UpdateHoldingRequest) (*EnrichedHolding, error) {
	h, err := s.repo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrHoldingNotFound
		}
		return nil, err
	}

	if req.Shares != nil {
		if *req.Shares <= 0 {
			return nil, errors.New("shares must be positive")
		}
		h.Shares = *req.Shares
	}
	if req.BuyPrice != nil {
		if *req.BuyPrice <= 0 {
			return nil, errors.New("buy_price must be positive")
		}
		h.BuyPrice = *req.BuyPrice
	}
	if req.BuyDate != nil {
		bd, err := parseBuyDate(*req.BuyDate)
		if err != nil {
			return nil, err
		}
		h.BuyDate = bd
	}
	if req.Notes != nil {
		h.Notes = strings.TrimSpace(*req.Notes)
	}

	if err := s.repo.Update(h); err != nil {
		return nil, err
	}

	e := s.enrichOne(*h, nil)
	return &e, nil
}

func (s *PortfolioService) Delete(userID, id string) error {
	err := s.repo.Delete(id, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repositories.ErrHoldingNotFound
	}
	return err
}

func parseBuyDate(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, errors.New("buy_date must be YYYY-MM-DD")
	}
	return &t, nil
}
