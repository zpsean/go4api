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
            MCharR13, 
            MCharR14, 
            MCharR15,
            MCharR16,
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
            MFloatR1, 
            MFloatR2, 
            MFloatR3,
            MFloatR4, 
            MFloatR5, 
            MFloatR6,
            MFloatR7,
        },
        "MBool": []interface{} {
            MBoolR1, 
            MBoolR2, 
            MBoolR3,
        },
        "MArray": []interface{} {
            MArrayR1,
            MArrayR2,
            MArrayR3,
            MArrayR4,
            MArrayR5,
            MArrayR6,
            MArrayR7,
            MArrayR8,
            MArrayR9,
            MArrayR10,
            MArrayR11,
            MArrayR12,
            MArrayR13,
            MArrayR14,
            MArrayR15,
            MArrayR16,
            MArrayR17,
            MArrayR18,
            MArrayR19,
        },
    }

    return RulesMapping[key]
}