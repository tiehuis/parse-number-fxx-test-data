// Copyright 2020 The ParseNumberFxxTestData Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slowstrconv_test

import (
	"fmt"
	"log"
	"math/big"

	"github.com/nigeltao/parse-number-fxx-test-data/slowstrconv"
)

func ExampleParseFloatFromBytes() {
	inputs := []string{
		"0.3",
		"3.141592653589793",
		"640",
	}

	for i, input := range inputs {
		res, err := slowstrconv.ParseFloatFromBytes([]byte(input))
		if err != nil {
			log.Fatal(err)
		}

		if i != 0 {
			fmt.Println()
		}
		fmt.Printf("Parsing %q\n", input)

		{
			fmt.Printf("==== F16: 0x%04X\n", res.F16)
			fmt.Printf("%s\n", input)
			man := int64((res.F16 & 0x03FF) | 0x0400)
			exp := int(res.F16>>10) - 15 - 10
			e := big.NewInt(0).Lsh(big.NewInt(1), uint(-exp))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man-1), e).FloatString(16))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man+0), e).FloatString(16))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man+1), e).FloatString(16))
		}

		{
			fmt.Printf("==== F32: 0x%08X\n", res.F32)
			fmt.Printf("%s\n", input)
			man := int64((res.F32 & 0x007F_FFFF) | 0x0080_0000)
			exp := int(res.F32>>23) - 127 - 23
			e := big.NewInt(0).Lsh(big.NewInt(1), uint(-exp))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man-1), e).FloatString(32))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man+0), e).FloatString(32))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man+1), e).FloatString(32))
		}

		{
			fmt.Printf("==== F64: 0x%016X\n", res.F64)
			fmt.Printf("%s\n", input)
			man := int64((res.F64 & 0x000F_FFFF_FFFF_FFFF) | 0x0010_0000_0000_0000)
			exp := int(res.F64>>52) - 1023 - 52
			e := big.NewInt(0).Lsh(big.NewInt(1), uint(-exp))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man-1), e).FloatString(64))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man+0), e).FloatString(64))
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(big.NewInt(man+1), e).FloatString(64))
		}

		{
			fmt.Printf("==== F128: 0x%016X%016X\n", res.F128.Hi, res.F128.Lo)
			fmt.Printf("%s\n", input)

			man := big.NewInt(int64((res.F128.Hi & 0x0000_FFFF_FFFF_FFFF) | 0x0001_0000_0000_0000))
			man.Lsh(man, 63)
			man.Add(man, big.NewInt(int64(res.F128.Lo>>1)))
			man.Lsh(man, 1)
			man.Add(man, big.NewInt(int64(res.F128.Lo&1)))

			exp := int(res.F128.Hi>>48) - 16383 - 112
			e := big.NewInt(0).Lsh(big.NewInt(1), uint(-exp))

			man.Sub(man, big.NewInt(1)) // man-1
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(man, e).FloatString(128))

			man.Add(man, big.NewInt(1)) // man+0
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(man, e).FloatString(128))

			man.Add(man, big.NewInt(1)) // man+1
			fmt.Printf("%v\n", big.NewRat(0, 1).SetFrac(man, e).FloatString(128))
		}
	}

	// Output:
	// Parsing "0.3"
	// ==== F16: 0x34CD
	// 0.3
	// 0.2998046875000000
	// 0.3000488281250000
	// 0.3002929687500000
	// ==== F32: 0x3E99999A
	// 0.3
	// 0.29999998211860656738281250000000
	// 0.30000001192092895507812500000000
	// 0.30000004172325134277343750000000
	// ==== F64: 0x3FD3333333333333
	// 0.3
	// 0.2999999999999999333866185224906075745820999145507812500000000000
	// 0.2999999999999999888977697537484345957636833190917968750000000000
	// 0.3000000000000000444089209850062616169452667236328125000000000000
	// ==== F128: 0x3FFD3333333333333333333333333333
	// 0.3
	// 0.29999999999999999999999999999999994222210166838292440832066172245218044385695055353835414280183613300323486328125000000000000000
	// 0.29999999999999999999999999999999999037035027806382073472011028707536340730949175892305902380030602216720581054687500000000000000
	// 0.30000000000000000000000000000000003851859888774471706111955885169854637076203296430776390479877591133117675781250000000000000000
	//
	// Parsing "3.141592653589793"
	// ==== F16: 0x4248
	// 3.141592653589793
	// 3.1386718750000000
	// 3.1406250000000000
	// 3.1425781250000000
	// ==== F32: 0x40490FDB
	// 3.141592653589793
	// 3.14159250259399414062500000000000
	// 3.14159274101257324218750000000000
	// 3.14159297943115234375000000000000
	// ==== F64: 0x400921FB54442D18
	// 3.141592653589793
	// 3.1415926535897926719087536184815689921379089355468750000000000000
	// 3.1415926535897931159979634685441851615905761718750000000000000000
	// 3.1415926535897935600871733186068013310432434082031250000000000000
	// ==== F128: 0x4000921FB54442D17BD21B8D78573DE5
	// 3.141592653589793
	// 3.14159265358979299999999999999999964658137890098464069769572079468241201927192207676853286102414131164550781250000000000000000000
	// 3.14159265358979300000000000000000003176736777843181130889130931166787572689225171984617190901190042495727539062500000000000000000
	// 3.14159265358979300000000000000000041695335665587898192008689782865333943451258136292381095699965953826904296875000000000000000000
	//
	// Parsing "640"
	// ==== F16: 0x6100
	// 640
	// 639.5000000000000000
	// 640.0000000000000000
	// 640.5000000000000000
	// ==== F32: 0x44200000
	// 640
	// 639.99993896484375000000000000000000
	// 640.00000000000000000000000000000000
	// 640.00006103515625000000000000000000
	// ==== F64: 0x4084000000000000
	// 640
	// 639.9999999999998863131622783839702606201171875000000000000000000000
	// 640.0000000000000000000000000000000000000000000000000000000000000000
	// 640.0000000000001136868377216160297393798828125000000000000000000000
	// ==== F128: 0x40084000000000000000000000000000
	// 640
	// 639.99999999999999999999999999999990139238684737352432353392933965172129084919561137212440371513366699218750000000000000000000000000
	// 640.00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	// 640.00000000000000000000000000000009860761315262647567646607066034827870915080438862787559628486633300781250000000000000000000000000
}
