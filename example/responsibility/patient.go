package main

// Patient 病人、患者
type Patient struct {
	name              string
	registrationDone  bool // 登记完成
	doctorCheckUpDone bool
	medicineDone      bool
	paymentDone       bool
}
