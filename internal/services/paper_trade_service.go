package services

import (
	"errors"
	"strings"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"github.com/google/uuid"
)

type PaperTradeInput struct {
	TickerSymbol    string     `json:"ticker_symbol"`
	Side            string     `json:"side"`
	Status          string     `json:"status"`
	Source          string     `json:"source"`
	Quantity        float64    `json:"quantity"`
	EntryPrice      float64    `json:"entry_price"`
	ExitPrice       *float64   `json:"exit_price"`
	StopLoss        *float64   `json:"stop_loss"`
	TakeProfit      *float64   `json:"take_profit"`
	TrailingStopPct *float64   `json:"trailing_stop_pct"`
	Notes           *string    `json:"notes"`
	OpenedAt        *time.Time `json:"opened_at"`
	ClosedAt        *time.Time `json:"closed_at"`
}

type PaperTradeService struct {
	repo *repositories.PaperTradeRepository
}

func NewPaperTradeService(repo *repositories.PaperTradeRepository) *PaperTradeService {
	return &PaperTradeService{repo: repo}
}

func (s *PaperTradeService) List(userID string) ([]models.PaperTrade, error) {
	return s.repo.ListByUser(userID, 500)
}

func (s *PaperTradeService) Create(userID string, in PaperTradeInput) (*models.PaperTrade, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	ticker := strings.ToUpper(strings.TrimSpace(in.TickerSymbol))
	if ticker == "" {
		return nil, errors.New("ticker_symbol is required")
	}
	if in.Quantity <= 0 || in.EntryPrice <= 0 {
		return nil, errors.New("quantity and entry_price must be positive")
	}
	side := orDefault(in.Side, "buy")
	status := orDefault(in.Status, "planned")
	source := orDefault(in.Source, "plan")
	if err := validatePaperTradeEnums(side, status, source); err != nil {
		return nil, err
	}
	trade := &models.PaperTrade{
		UserID:          uid,
		TickerSymbol:    ticker,
		Side:            side,
		Status:          status,
		Source:          source,
		Quantity:        in.Quantity,
		EntryPrice:      in.EntryPrice,
		ExitPrice:       in.ExitPrice,
		StopLoss:        in.StopLoss,
		TakeProfit:      in.TakeProfit,
		TrailingStopPct: in.TrailingStopPct,
		ClosedAt:        in.ClosedAt,
		OpenedAt:        time.Now(),
	}
	if in.Notes != nil {
		trade.Notes = *in.Notes
	}
	if in.OpenedAt != nil {
		trade.OpenedAt = *in.OpenedAt
	}
	trade.PnL = computePaperPnL(trade)
	if err := s.repo.Create(trade); err != nil {
		return nil, err
	}
	return trade, nil
}

// Update applies a partial update. Closing a trade (exit_price set) recomputes P&L.
func (s *PaperTradeService) Update(userID, id string, in PaperTradeInput) (*models.PaperTrade, error) {
	updates := map[string]interface{}{}
	if in.Status != "" {
		if err := validatePaperTradeEnums("buy", in.Status, "plan"); err != nil {
			return nil, err
		}
		updates["status"] = in.Status
	}
	if in.ExitPrice != nil {
		updates["exit_price"] = *in.ExitPrice
	}
	if in.StopLoss != nil {
		updates["stop_loss"] = *in.StopLoss
	}
	if in.TakeProfit != nil {
		updates["take_profit"] = *in.TakeProfit
	}
	if in.TrailingStopPct != nil {
		updates["trailing_stop_pct"] = *in.TrailingStopPct
	}
	if in.Notes != nil {
		updates["notes"] = *in.Notes
	}
	if in.ClosedAt != nil {
		updates["closed_at"] = *in.ClosedAt
	}
	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	if in.ExitPrice != nil {
		existing, err := s.repo.ListByUser(userID, 500)
		if err != nil {
			return nil, err
		}
		for i := range existing {
			if existing[i].ID.String() == id {
				merged := existing[i]
				merged.ExitPrice = in.ExitPrice
				updates["pnl"] = *computePaperPnL(&merged)
				break
			}
		}
	}
	return s.repo.Update(userID, id, updates)
}

func (s *PaperTradeService) Delete(userID, id string) error {
	return s.repo.Delete(userID, id)
}

func computePaperPnL(t *models.PaperTrade) *float64 {
	if t.ExitPrice == nil {
		return nil
	}
	pnl := (*t.ExitPrice - t.EntryPrice) * t.Quantity
	if t.Side == "sell" {
		pnl = -pnl
	}
	return &pnl
}

func orDefault(v, d string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return d
	}
	return v
}

func validatePaperTradeEnums(side, status, source string) error {
	switch side {
	case "buy", "sell":
	default:
		return errors.New("side must be buy or sell")
	}
	switch status {
	case "planned", "open", "closed", "cancelled":
	default:
		return errors.New("invalid status")
	}
	switch source {
	case "plan", "replay", "backtest", "neway":
	default:
		return errors.New("invalid source")
	}
	return nil
}
