/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/drivers"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
)

func main() {
	driverManager := plc4go.NewPlcDriverManager()
	//drivers.RegisterModbusDriver(driverManager)
	drivers.RegisterS7Driver(driverManager)

	// Get a connection to a remote PLC
	crc := driverManager.GetConnection("s7://192.168.216.121?remote-rack=0&remote-slot=1&controller-type=S7_1500")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		fmt.Printf("error connecting to PLC: %s", connectionResult.Err.Error())
		return
	}
	connection := connectionResult.Connection

	// Make sure the connection is closed at the end
	defer connection.BlockingClose()

	// Prepare a read-request
	readRequest, err := connection.ReadRequestBuilder().
		AddQuery("field_Q0_0", "%Q0.0:BOOL").
		AddQuery("field_I0_0", "%I0.0:BOOL").
		AddQuery("field_QW2", "%M0:INT").
		Build()
	if err != nil {
		fmt.Printf("error preparing read-request: %s", connectionResult.Err.Error())
		return
	}

	// Execute a read-request
	rrc := readRequest.Execute()

	// Wait for the response to finish
	rrr := <-rrc
	if rrr.Err != nil {
		fmt.Printf("error executing read-request: %s", rrr.Err.Error())
		return
	}

	// Do something with the response
	if rrr.Response.GetResponseCode("field_I0_0") != model.PlcResponseCode_OK {
		fmt.Printf("error an non-ok return code: %s", rrr.Response.GetResponseCode("field").GetName())
		return
	}

	value := rrr.Response.GetValue("field_Q0_0")
	fmt.Printf("Got result %t\n", value.GetBool())
	value = rrr.Response.GetValue("field_I0_0")
	fmt.Printf("Got result %v\n", value)
	value = rrr.Response.GetValue("field_QW2")
	fmt.Printf("Got result %d\n", value.GetInt16())

}
