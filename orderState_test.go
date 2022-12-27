package stateless_test

import (
	"context"
	"fmt"
	"github.com/qmuntal/stateless"
	"testing"
)

const (
	triggerCheckoutPaid            = "checkoutPaid"
	triggerDeadlinePassed          = "deadlinePassed"
	triggerUserConfirmDelivery     = "userConfirmDelivery"
	triggerCancelCheckOut          = "cancelCheckOut"
	triggerSysCheckOutNoMatch      = "sysCheckOutNoMatch"
	triggerCheckFailed             = "checkFailed"
	triggerBuyerCancelOrder        = "buyerCancelOrder"
	triggerSellerCancelOrder       = "sellerCancelOrder"
	triggerSellerAcceptOrderCancel = "sellerAcceptCancelOrder"
	triggerSellerRejectOrderCancel = "sellerRejectCancelOrder"
	triggerRefundPaid              = "refundPaid"
	triggerReturnRequest           = "returnRequested"
	triggerReturnRequestCancel     = "returnRequestCancel"
	triggerReturnAccepted          = "returnAccepted"
)

const (
	stateUnpaid           = "Unpaid"
	statePaid             = "Paid"
	stateCompleted        = "Completed"
	stateCancelPending    = "CancelPending"
	stateCancelProcessing = "CancelProcessing"
	stateCancelCompleted  = "CancelCompleted"
	stateReturnProcessing = "ReturnProcessing"
	stateReturnCompleted  = "ReturnCompleted"
	stateInvalid          = "Invalid"
)

func TestOrderState(t *testing.T) {
	orderStateMachine := stateless.NewStateMachine(stateUnpaid)

	orderStateMachine.Configure(stateUnpaid).
		Permit(triggerCheckoutPaid, statePaid, CheckOutStatePaid).
		Permit(triggerCancelCheckOut, stateInvalid, UserCancelCheckout).
		Permit(triggerSysCheckOutNoMatch, stateInvalid, SysUpdateCheckoutUnMatch).
		Permit(triggerCheckFailed, stateInvalid, CheckOutFailed).
		Permit(triggerBuyerCancelOrder, stateCancelPending, BuyerCancelCOD)
	orderStateMachine.Configure(statePaid).
		OnEntryFrom(triggerDeadlinePassed, ConfirmDeliveryDeadlinePass).
		InternalTransition(triggerUserConfirmDelivery, UserConfirmDeliveryAction).
		Permit(triggerDeadlinePassed, stateCompleted).
		Permit(triggerUserConfirmDelivery, stateCompleted).
		Permit(triggerBuyerCancelOrder, stateCancelPending, BuyerCancelNonCodOrder).
		Permit(triggerSellerCancelOrder, stateCancelProcessing, SellerCancelOrder).
		Permit(triggerReturnRequest, stateReturnProcessing)

	orderStateMachine.Configure(stateCancelPending).
		Permit(triggerSellerAcceptOrderCancel, stateCancelProcessing, SellerAcceptNonCodCancel).
		Permit(triggerSellerRejectOrderCancel, stateUnpaid, SellerRejectCodCancel).
		Permit(triggerSellerRejectOrderCancel, statePaid, SellerRejectNonCodCancel).
		Permit(triggerSellerAcceptOrderCancel, stateInvalid, SellerAcceptCodCancel)
	orderStateMachine.Configure(stateCancelProcessing).
		Permit(triggerRefundPaid, stateCancelCompleted)

	orderStateMachine.Configure(stateReturnProcessing).
		Permit(triggerReturnRequestCancel, statePaid)
	orderStateMachine.Configure(stateReturnProcessing).
		Permit(triggerReturnAccepted, stateReturnCompleted)

	grahStr := orderStateMachine.ToGraph()
	fmt.Println(grahStr)

	canPaidOverPaid, _ := orderStateMachine.CanFire(triggerCheckoutPaid, "overpaid")
	fmt.Println(canPaidOverPaid)
	canPaidPending, _ := orderStateMachine.CanFire(triggerCheckoutPaid, "pending")
	fmt.Println(canPaidPending)
	orderStateMachine.Fire(triggerCheckoutPaid, "paid")
	currentState := orderStateMachine.MustState()
	fmt.Printf("current state is: %v\n", currentState)
	canTransitionTo, _ := orderStateMachine.IsInState(stateUnpaid)
	fmt.Printf("can transition to state %v, result:%v \n", stateUnpaid, canTransitionTo)

}

func CheckOutStatePaid(_ context.Context, args ...interface{}) bool {
	//condition guards: checkout state is paid/overpaid
	return args[0].(string) == "paid" || args[0].(string) == "overpaid"
}

func UserCancelCheckout(_ context.Context, args ...interface{}) bool {
	//condition guards: user cancel && order state == unpaid
	return args[0].(string) == "User" && args[1].(string) == "unpaid"
}

func SysUpdateCheckoutUnMatch(_ context.Context, args ...interface{}) bool {
	//condition guards: order state == unpaid && check state == noMatch
	return args[0].(string) == "unpaid" && args[1].(string) == "noMatch"
}

func CheckOutFailed(_ context.Context, args ...interface{}) bool {
	//condition guards: order state == unpaid && check state == failed
	return args[0].(string) == "unpaid" && args[1].(string) == "failed"
}

func BuyerCancelCOD(_ context.Context, args ...interface{}) bool {
	//condition guards: buyer trigger && checkout state == matching && order type== COD
	return args[0].(string) == "Buyer" && args[1].(string) == "matching" && args[2].(string) == "COD"
}

func ConfirmDeliveryDeadlinePass(_ context.Context, args ...interface{}) error {
	fmt.Printf("trigger%v order compeleted\n", triggerDeadlinePassed)
	return nil
}

func UserConfirmDeliveryAction(_ context.Context, _ ...interface{}) error {
	fmt.Printf("trigger%v order compeleted\n", triggerUserConfirmDelivery)
	return nil
}

func BuyerCancelNonCodOrder(_ context.Context, args ...interface{}) bool {
	//condition guards: checkout state == matching && order type== COD
	return args[0].(string) == "matching" && args[1].(string) != "COD"
}

func SellerCancelOrder(_ context.Context, args ...interface{}) bool {
	//condition guards: seller trigger
	return args[0].(string) == "Seller" && args[1].(string) == "matching"
}

func SellerAcceptNonCodCancel(_ context.Context, args ...interface{}) bool {
	//condition guards: seller accept cancel for non-COD orders
	return args[0].(string) == "Seller" && args[1].(string) != "COD"
}
func SellerAcceptCodCancel(_ context.Context, args ...interface{}) bool {
	//condition guards: seller reject cancel orders && for intially unpaid COD orders
	return args[0].(string) == "Seller" && args[1].(string) == "COD"
}

func SellerRejectCodCancel(_ context.Context, args ...interface{}) bool {
	//condition guards: seller reject cancel orders && for intially unpaid COD orders
	return args[0].(string) == "Seller" && args[1].(string) == "COD"
}

func SellerRejectNonCodCancel(_ context.Context, args ...interface{}) bool {
	//condition guards: seller reject cancel orders && for initally unpaid non COD orders
	return args[0].(string) == "Seller" && args[1].(string) != "COD"
}
