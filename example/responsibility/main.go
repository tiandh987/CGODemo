package main

func main() {

	cashier := &Cashier{}

	// Set next for medical department
	// 为医疗部门设置下一个
	medical := &Medical{}
	medical.setNext(cashier)

	// Set next for doctor department
	// 为医生科室设置下一个
	doctor := &Doctor{}
	doctor.setNext(medical)

	//Set next for reception department
	// 为接待部门设置下一个
	reception := &Reception{}
	reception.setNext(doctor)

	patient := &Patient{name: "abc"}
	// Patient visiting
	// 患者就诊
	reception.execute(patient)
}
