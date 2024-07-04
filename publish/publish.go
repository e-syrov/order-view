package main

import (
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	data := [][]byte{
		[]byte("{\n  \"order_uid\": \"b563feb7b2b84b6test\"}"),
	}

	/*data := [][]byte{
		[]byte("{\n  \"order_uid\": \"b563feb7b2b84b6test\",\n  \"track_number\": \"WBILMTESTTRACK\",\n  \"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Test Testov\",\n    \"phone\": \"+9720000000\",\n    \"zip\": \"2639809\",\n    \"city\": \"Kiryat Mozkin\",\n    \"address\": \"Ploshad Mira 15\",\n    \"region\": \"Kraiot\",\n    \"email\": \"test@gmail.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"b563feb7b2b84b6test\",\n    \"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"wbpay\",\n    \"amount\": 1817,\n    \"payment_dt\": 1637907727,\n    \"bank\": \"alpha\",\n    \"delivery_cost\": 1500,\n    \"goods_total\": 317,\n    \"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 9934930,\n      \"track_number\": \"WBILMTESTTRACK\",\n      \"price\": 453,\n      \"rid\": \"ab4219087a764ae0btest\",\n      \"name\": \"Mascaras\",\n      \"sale\": 30,\n      \"size\": \"0\",\n      \"total_price\": 317,\n      \"nm_id\": 2389212,\n      \"brand\": \"Vivienne Sabo\",\n      \"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"test\",\n  \"delivery_service\": \"meest\",\n  \"shardkey\": \"9\",\n  \"sm_id\": 99,\n  \"date_created\": \"2021-11-26T06:22:19Z\",\n  \"oof_shard\": \"1\"\n}"),
		[]byte("{\n  \"order_uid\": \"a123fdb7b2b84b6test1\",\n  \"track_number\": \"TRACK123\",\n  \"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Alice Brown\",\n    \"phone\": \"+1234567890\",\n    \"zip\": \"1234567\",\n    \"city\": \"New York\",\n    \"address\": \"5th Avenue 10\",\n    \"region\": \"NY\",\n    \"email\": \"alice@example.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"a123fdb7b2b84b6test1\",\n    \"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"paypal\",\n    \"amount\": 2000,\n    \"payment_dt\": 1637907727,\n    \"bank\": \"bank of america\",\n    \"delivery_cost\": 200,\n    \"goods_total\": 1800,\n    \"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 1234567,\n      \"track_number\": \"TRACK123\",\n      \"price\": 600,\n      \"rid\": \"ab1234567a764ae0btest\",\n      \"name\": \"Lipstick\",\n      \"sale\": 20,\n      \"size\": \"0\",\n      \"total_price\": 480,\n      \"nm_id\": 3456789,\n      \"brand\": \"MAC\",\n      \"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"customer1\",\n  \"delivery_service\": \"dhl\",\n  \"shardkey\": \"1\",\n  \"sm_id\": 101,\n  \"date_created\": \"2021-12-01T10:00:00Z\",\n  \"oof_shard\": \"1\"\n}"),
		[]byte("{\n  \"order_uid\": \"b456feb7b2b84b6test2\",\n  \"track_number\": \"TRACK456\",\n  \"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Bob Smith\",\n    \"phone\": \"+9876543210\",\n    \"zip\": \"7654321\",\n    \"city\": \"Los Angeles\",\n    \"address\": \"Sunset Blvd 100\",\n    \"region\": \"CA\",\n    \"email\": \"bob@example.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"b456feb7b2b84b6test2\",\n    \"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"stripe\",\n    \"amount\": 1500,\n    \"payment_dt\": 1637907727,\n    \"bank\": \"chase\",\n    \"delivery_cost\": 100,\n    \"goods_total\": 1400,\n    \"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 7654321,\n      \"track_number\": \"TRACK456\",\n      \"price\": 300,\n      \"rid\": \"cd1234567a764ae0btest\",\n      \"name\": \"Perfume\",\n      \"sale\": 10,\n      \"size\": \"0\",\n      \"total_price\": 270,\n      \"nm_id\": 9876543,\n      \"brand\": \"Chanel\",\n      \"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"customer2\",\n  \"delivery_service\": \"fedex\",\n  \"shardkey\": \"2\",\n  \"sm_id\": 102,\n  \"date_created\": \"2021-12-02T11:00:00Z\",\n  \"oof_shard\": \"2\"\n}"),
		[]byte("{\n  \"order_uid\": \"c789feb7b2b84b6test3\",\n  \"track_number\": \"TRACK789\",\n  \"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Charlie Johnson\",\n    \"phone\": \"+1112223333\",\n    \"zip\": \"1112223\",\n    \"city\": \"Chicago\",\n    \"address\": \"Michigan Ave 20\",\n    \"region\": \"IL\",\n    \"email\": \"charlie@example.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"c789feb7b2b84b6test3\",\n    \"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"square\",\n    \"amount\": 2500,\n    \"payment_dt\": 1637907727,\n    \"bank\": \"wells fargo\",\n    \"delivery_cost\": 300,\n    \"goods_total\": 2200,\n    \"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 1237890,\n      \"track_number\": \"TRACK789\",\n      \"price\": 500,\n      \"rid\": \"ef1234567a764ae0btest\",\n      \"name\": \"Watch\",\n      \"sale\": 15,\n      \"size\": \"0\",\n      \"total_price\": 425,\n      \"nm_id\": 4567890,\n      \"brand\": \"Rolex\",\n      \"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"customer3\",\n  \"delivery_service\": \"ups\",\n  \"shardkey\": \"3\",\n  \"sm_id\": 103,\n  \"date_created\": \"2021-12-03T12:00:00Z\",\n  \"oof_shard\": \"3\"\n}"),
		[]byte("{\n  \"order_uid\": \"d012feb7b2b84b6test4\",\n  \"track_number\": \"TRACK012\",\n  \"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Dana White\",\n    \"phone\": \"+4445556666\",\n    \"zip\": \"4445556\",\n    \"city\": \"Houston\",\n    \"address\": \"Main St 30\",\n    \"region\": \"TX\",\n    \"email\": \"dana@example.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"d012feb7b2b84b6test4\",\n    \"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"authorize.net\",\n    \"amount\": 3000,\n    \"payment_dt\": 1637907727,\n    \"bank\": \"citibank\",\n    \"delivery_cost\": 400,\n    \"goods_total\": 2600,\n    \"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 4560123,\n      \"track_number\": \"TRACK012\",\n      \"price\": 700,\n      \"rid\": \"gh1234567a764ae0btest\",\n      \"name\": \"Shoes\",\n      \"sale\": 25,\n      \"size\": \"0\",\n      \"total_price\": 525,\n      \"nm_id\": 6789012,\n      \"brand\": \"Nike\",\n      \"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"customer4\",\n  \"delivery_service\": \"usps\",\n  \"shardkey\": \"4\",\n  \"sm_id\": 104,\n  \"date_created\": \"2021-12-04T13:00:00Z\",\n  \"oof_shard\": \"4\"\n}"),
	}*/

	conn, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatalf("Ошибка отправления в канал")
	}
	defer conn.Close()
	for _, d := range data {
		publish(conn, "orders-channel", d)
	}

}
func publish(conn stan.Conn, nameCh string, message []byte) {
	err := conn.Publish(nameCh, message)
	if err != nil {
		log.Fatalf("Ошибка публикации: %v", err)
	}
}
