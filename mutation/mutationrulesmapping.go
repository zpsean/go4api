/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package mutation

import (                                                                                                                                             

)


func MutationRulesMapping(key string) []interface{} {
    //
    RulesMapping := map[string][]interface{} {
        "MChar": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Common_Set_To_Null,
            M_Common_Set_To_Single_Plus_Sign,
            M_Common_Set_To_Single_Percent_Sign,
            M_Common_Set_To_Single_Number_Sign,
            M_Common_Set_To_Single_Minus_Sign,
            M_Common_Set_To_Single_Exclamation_Sign,
            M_Common_Set_To_Single_e,
            M_Common_Set_To_Single_E,
            M_Common_Set_To_Single_Dot_Sign,
            M_Common_Set_To_Single_Dollar_Sign,
            M_Common_Set_To_Single_Caret_Sign,
            M_Common_Set_To_Single_At_Sign,
            M_Common_Set_To_Single_Asterisk_Sign,
            M_Common_Set_To_Single_Ampersand_Sign,
            M_Char_Add_Prefix_One_Blank,
            M_Char_Replace_All_Blank,
            M_Char_Replace_Prefix_One_Blank,
            M_Char_Replace_Prefix_None_ASCII,
            M_Char_Add_Suffix_One_Blank, 
            M_Char_Replace_Suffix_One_Blank,
            M_Char_Add_Mid_One_Blank, 
            M_Char_Add_Mid_None_ASCII,
            M_Char_Replace_Mid_One_Blank, 
            M_Char_Replace_Mid_One_E,
            M_Char_Replace_Mid_One_Negative_Sign,
            M_Char_Replace_Mid_None_ASCII,
            M_Char_Set_To_One_Char,
            M_Char_Repeat_50_Times, 
            M_Char_Replace_Prefix_Percentage, 
            M_Char_Replace_Prefix_Point,
            M_Char_Replace_Prefix_Caret,
            M_Char_Replace_Prefix_Dollar,
            M_Char_Replace_Prefix_Star,
            M_Char_Set_To_Int, 
            M_Char_Set_To_Float,
            M_Common_Set_To_Bool_True,
            M_Common_Set_To_Bool_False,
            M_Char_Set_To_Array,
        },
        "MCharNumeric": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MCharAlpha": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MCharAlphaNumeric": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MCharTime": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MCharEmail": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MCharIp": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MInt": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank,
            M_Int_Set_To_Zero,
            M_Int_Set_To_One,
            M_Int_Set_To_Negative_One,
            M_Int_Set_To_MaxInt32,
            M_Int_Set_To_MinInt32,
            M_Int_Set_To_MaxInt64,
            M_Int_Set_To_MinInt64,
        },
        "MIntTime": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank, 
            M_Char_Add_Prefix_One_Blank,
        },
        "MFloat": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank,
            M_Float_Set_To_E,
            M_Float_Set_To_Postive_Float,
            M_Float_Set_To_Negative_Float,
            M_Float_Set_To_MaxFloat32,
            M_Float_Set_To_SmallestNonzeroFloat32,
            M_Float_Set_To_MaxFloat64,
            M_Float_Set_To_SmallestNonzeroFloat64,
        },
        "MBool": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank,
            M_Common_Set_To_Null,
            M_Common_Set_To_Bool_True,
            M_Common_Set_To_Bool_False,
            M_Bool_Set_To_Zero,
        },
        "MArray": []interface{} {
            M_Common_Set_To_Empty, 
            M_Common_Set_To_One_Blank,
            M_Array_Set_To_Empty_Array,
            M_Array_Remove_One_Item_Random,
            M_Array_Set_Only_One_Item,
            M_Array_Duplicate_One_Item_Random,
            M_Array_Append_Another_Type_Item,
            M_Array_Replace_Another_Type_Item,
            M_Array_Replace_One_Item_Null,
            M_Array_Replace_One_Item_Bool_True,
            M_Array_Replace_One_Item_Bool_False,
            M_Array_Set_To_Only_One_Null,
            M_Array_Set_To_Only_One_Int,
            M_Array_Set_To_Only_One_String,
            M_Array_Set_To_Only_One_Bool_True,
            M_Array_Set_To_Only_One_Bool_False,
            M_Array_Set_To_Int,
            M_Array_Set_To_String,
            M_Common_Set_To_Bool_True,
            M_Common_Set_To_Bool_False,
            M_Common_Set_To_Null,
        },
        "MMap": []interface{} {
            M_Map_Set_To_Empty_Map,
        },
    }

    return RulesMapping[key]
}