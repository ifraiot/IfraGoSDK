# IfraGoSDK
## Installation
```
go get github.com/ifraiot/IfraGoSDK
```
## Using Example
```
func main() {

	topic := flag.String("topic", "", "Organization's topic")
	username := flag.String("username", "", "MQTT username")
	password := flag.String("password", "", "MQTT password")
	flag.Parse()

	ifraCon := ifraGoSDK.NewIFRA(
		*topic,
		*username,
		*password)

	for {
		memory, err := memory.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

 
		percent, _ := cpu.Percent(time.Second, true)
		//Memory
		ifraCon.AddMeasurement("memory_total", float64(memory.Total)/1024/1024)
		ifraCon.AddMeasurement("memory_used", float64(memory.Used)/1024/1024)
		ifraCon.AddMeasurement("memory_cached", float64(memory.Cached)/1024/1024)
		ifraCon.AddMeasurement("memory_free", float64(memory.Free)/1024/1024)

		//CPU
		ifraCon.AddMeasurement("cpu_sys", percent[cpu.CPSys])
		ifraCon.AddMeasurement("cpu_user", percent[cpu.CPUser])
		ifraCon.AddMeasurement("cpu_usage", percent[cpu.CPUser]+percent[cpu.CPSys])
		ifraCon.AddMeasurement("cpu_idle", percent[cpu.CPIdle])
		ifraCon.Send()

		//battery
		batteries, err := battery.GetAll()
		if err != nil {
			fmt.Println("Could not get battery info!")
			return
		}
		for _, battery := range batteries {

			if battery.State.String() == "Full" {
				ifraCon.AddMeasurement("battery_state", 1)
			} else {
				ifraCon.AddMeasurement("battery_state", 0)
			}

			ifraCon.AddMeasurement("battery_capacity", battery.Current)
			ifraCon.AddMeasurement("battery_last_capacity", battery.Full)
			ifraCon.AddMeasurement("battery_charge_rate", battery.Current/battery.Full*100)
		}

		time.Sleep(5 * time.Second)
	}
}
```
