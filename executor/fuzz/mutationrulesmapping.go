/*
 * go4api - a api testing tool written in Go
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


func MutationRulesMapping(key string) []interface{} {
    //
    RulesMapping := map[string][]interface{} {
        "MChar": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
            MCharR4, 
            MCharR5, 
            MCharR6,
            MCharR7, 
            MCharR8, 
            MCharR9,
            MCharR10, 
            MCharR11, 
            MCharR12,
            // MCharR13, 
            // MCharR14, 
            // MCharR15,
            // MCharR16,
        },
        "MCharNumeric": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MCharAlpha": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MCharAlphaNumeric": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MCharTime": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MCharEmail": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MCharIp": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MInt": []interface{} {
            MIntR1, 
            MIntR2,
            MIntR3,
            MIntR4,
            MIntR5,
            MIntR6,
        },
        "MIntTime": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MFloat": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MBool": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
        "MArray": []interface{} {
            MCharR1, 
            MCharR2, 
            MCharR3,
        },
    }

    return RulesMapping[key]
}