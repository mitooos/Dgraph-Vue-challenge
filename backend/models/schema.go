package models

const schema = `
	id: string @index(exact, term) .
	name: string .
	age: int .
	price: int .
	date: datetime @index(day) .
	ip: string @index(exact) .
	device: string .

	transactions: [uid] .
	buyer: uid .
	products: [uid] @count @reverse .

	type Buyer{
		id
		name
		age
		transactions
		date
	}

	type Product{
		id
		name
		price
		date
	}

	type Transaction{
		id
		date
		buyer
		ip
		device
		products
	}
`
