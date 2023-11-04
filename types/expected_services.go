package types

type EscrowService interface {
	NewContract() (uint64, error)
}
