// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package utils

const SenderTVC = "te6ccgECEQEAAhkAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAib/APSkICLAAZL0oOGK7VNYMPShCAQBCvSkIPShBQIDz0AHBgAnTtRNDT/9M/0wDRf/hh+Gb4Y/higAJV+ELIy//4Q88LP/hGzwsAye1UgCASAMCQEC/woB/n8h7UTQINdJwgGOENP/0z/TANF/+GH4Zvhj+GKOGPQFcAGAQPQO8r3XC//4YnD4Y3D4Zn/4YeLTAAGOEoECANcYIPkBWPhCIPhl+RDyqN7TPwGOHvhDIbkgnzAg+COBA+iogggbd0Cgud6S+GPggDTyNNjTHyHBAyKCEP////0LACC8sZLyPOAB8AH4R26S8jzeAgEgDg0Anb1Fqvn/wgt0ca9qJoEGuk4QDHCGn/6Z/pgGi//DD8M3wx/DFHDHoCuADAIHoHeV7rhf/8MTh8Mbh8Mz/8MPFvfCN5Obj8M2j8AHgAv/wzwCASAQDwDju/VJio+EFukvAC3vpA0gDXDX+V1NHQ03/f1w3/ldTR0NP/39cN/5XU0dDT/9/R+EL4RSBukjBw3rry4Gz4ACIlJcjPhYDKAHPPQM4B+gKAac9Az4HPg8jPkYZ8I24jzwv/Is8L/83JcfsAXwXwAX/4Z4AFzdcCLQ1wsDqTgA3CHHANwh0x8h3SHBAyKCEP////28sZLyPOAB8AH4R26S8jze"
const SenderABI = `{
	"ABI version": 2,
	"header": ["time"],
	"functions": [
		{
			"name": "sendData",
			"inputs": [
				{"name":"destination","type":"address"},
				{"name":"bounce","type":"bool"},
				{"name":"value","type":"uint128"},
				{"name":"data","type":"uint256"},
				{"name":"destinationChainId","type":"uint256"}
			],
			"outputs": [
			]
		},
		{
			"name": "constructor",
			"inputs": [
			],
			"outputs": [
			]
		}
	],
	"data": [
	],
	"events": [
	]
}`
