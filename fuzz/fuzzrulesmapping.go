/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package fuzz

import (

)


func FuzzRulesMapping(key string) []interface{} {
    //
    RulesMapping := map[string][]interface{} {
        "FCharValid": []interface{} {
            FCharValidR1, 
            FCharValidR2, 
            FCharValidR3,
        },
        "FCharInvalid": []interface{} {
            FCharValidR1, 
            FCharValidR2, 
            FCharValidR3,
        },
        "FCharNumericValid": []interface{} {
            FCharValidR1, 
            FCharValidR2, 
            FCharValidR3,
        },
        "FCharAlphaValid": []interface{} {
            FCharValidR1, 
            FCharValidR2, 
            FCharValidR3,
        },
        
    }

    return RulesMapping[key]
}