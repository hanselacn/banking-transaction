package consts

const (
	TxTypeDEBIT  = "super_admin"
	TxTypeCREDIT = "admin"

	TxActionWITHDRAWAL = "customer"
	TxActionDEPOSIT    = "customer"
	TxActionTRANSFER   = "customer"
	TxActionPURCHASE   = "customer"

	TxStatusINPROGRESS = "IN_PROGRESS"
	TxStatusCOMPLETED  = "COMPLETED"
	TxStatusFAILED     = "FAILED"
)
