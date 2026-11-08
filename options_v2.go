package qe_connector

import (
	"errors"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

// OptionMetadata is the server-validated contract identity saved with an order.
// It contains no live quote, Greeks, or target-position snapshot.
type OptionMetadata struct {
	OptionType          string    `json:"optionType"`
	Strike              string    `json:"strike"`
	ExpirationTimestamp FlexInt64 `json:"expirationTimestamp"`
	ContractSize        string    `json:"contractSize"`
	ContractType        string    `json:"contractType"`
	UnderlyingIndex     string    `json:"underlyingIndex"`
}

// OptionPriceStep applies when the premium is greater than or equal to AbovePrice.
type OptionPriceStep struct {
	AbovePrice string `json:"abovePrice"`
	TickSize   string `json:"tickSize"`
}

var optionDecimalPattern = regexp.MustCompile(`^[0-9]{1,20}(\.[0-9]{1,10})?$`)

var optionTriggerPricePattern = regexp.MustCompile(`^\+?([0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE]([+-]?[0-9]+))?$`)

func validOptionDecimal(value string, allowZero bool) bool {
	value = strings.TrimSpace(value)
	if !optionDecimalPattern.MatchString(value) {
		return false
	}
	n, ok := new(big.Rat).SetString(value)
	return ok && (n.Sign() > 0 || allowZero && n.Sign() == 0)
}

// Match backend normalizeTriggerPrice without rounding or expanding exponents.
// A valid zero/blank trigger is disabled; a positive trigger is enabled.
func optionTriggerEnabled(raw string) (bool, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, nil
	}
	parts := optionTriggerPricePattern.FindStringSubmatch(raw)
	if parts == nil {
		return false, errors.New("OPTION trigger price must fit non-negative DECIMAL(30,10)")
	}
	integer, fraction, _ := strings.Cut(parts[1], ".")
	allDigits := integer + fraction
	digits := strings.TrimLeft(allDigits, "0")
	if digits == "" {
		return false, nil
	}
	var exponent int64
	if parts[2] != "" {
		var err error
		exponent, err = strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return false, errors.New("OPTION trigger price must fit non-negative DECIMAL(30,10)")
		}
	}
	significant := strings.TrimRight(digits, "0")
	integerPlaces := int64(len(integer)) + exponent - int64(len(allDigits)-len(digits))
	if integerPlaces > 20 || int64(len(significant))-integerPlaces > 10 {
		return false, errors.New("OPTION trigger price must fit non-negative DECIMAL(30,10)")
	}
	return true, nil
}

// Only static OPTION request constraints belong here. Published instruments,
// quantity/price steps, expiry and live positions remain server-validated.
func (s *CreateMasterOrderV2Service) validateOption() error {
	if s.exchange != trading_enums.ExchangeDeribit {
		return errors.New("OPTION is supported only on Deribit")
	}
	if s.orderNotional != nil || s.isMargin != nil && *s.isMargin {
		return errors.New("OPTION does not support orderNotional or isMargin=true")
	}
	target := s.isTargetPosition != nil && *s.isTargetPosition
	if s.totalQuantity == nil || !validOptionDecimal(*s.totalQuantity, target) {
		return errors.New("OPTION requires a positive base-asset totalQuantity (zero is allowed for a target)")
	}
	if s.worstPrice == nil || !validOptionDecimal(*s.worstPrice, false) {
		return errors.New("OPTION requires a positive premium worstPrice")
	}
	if target && (s.expectedExecutionSide == nil || (*s.expectedExecutionSide != trading_enums.OrderSideBuy && *s.expectedExecutionSide != trading_enums.OrderSideSell)) {
		return errors.New("OPTION targets require expectedExecutionSide=buy or sell from the position preview")
	}
	for _, price := range []*string{s.triggerPrice, s.triggerPrice1} {
		if price == nil {
			continue
		}
		enabled, err := optionTriggerEnabled(*price)
		if err != nil {
			return err
		}
		if enabled && (s.maxTriggerWaitSecs == nil || *s.maxTriggerWaitSecs <= 0 || math.IsNaN(*s.maxTriggerWaitSecs) || math.IsInf(*s.maxTriggerWaitSecs, 0)) {
			return errors.New("OPTION price triggers require a finite positive maxTriggerWaitSecs")
		}
	}
	return nil
}
