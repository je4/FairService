package re3data

type RepositoryType string

const (
	RepositoryTypeDisciplinary      RepositoryType = "disciplinary"
	RepositoryTypeGovernmental      RepositoryType = "governmental"
	RepositoryTypeInstitutional     RepositoryType = "institutional"
	RepositoryTypeMultidisciplinary RepositoryType = "multidisciplinary"
	RepositoryTypeProjectRelated    RepositoryType = "project-related"
	RepositoryTypeOther             RepositoryType = "other"
)

type SubjectType string

const (
	SubjectType1       SubjectType = "1"
	SubjectType11      SubjectType = "11"
	SubjectType101     SubjectType = "101"
	SubjectType10101   SubjectType = "10101"
	SubjectType10102   SubjectType = "10102"
	SubjectType10103   SubjectType = "10103"
	SubjectType10104   SubjectType = "10104"
	SubjectType10105   SubjectType = "10105"
	SubjectType102     SubjectType = "102"
	SubjectType10201   SubjectType = "10201"
	SubjectType10202   SubjectType = "10202"
	SubjectType10203   SubjectType = "10203"
	SubjectType10204   SubjectType = "10204"
	SubjectType103     SubjectType = "103"
	SubjectType10301   SubjectType = "10301"
	SubjectType10302   SubjectType = "10302"
	SubjectType10303   SubjectType = "10303"
	SubjectType104     SubjectType = "104"
	SubjectType10401   SubjectType = "10401"
	SubjectType10402   SubjectType = "10402"
	SubjectType10403   SubjectType = "10403"
	SubjectType105     SubjectType = "105"
	SubjectType10501   SubjectType = "10501"
	SubjectType10502   SubjectType = "10502"
	SubjectType10503   SubjectType = "10503"
	SubjectType10504   SubjectType = "10504"
	SubjectType106     SubjectType = "106"
	SubjectType10601   SubjectType = "10601"
	SubjectType10602   SubjectType = "10602"
	SubjectType10603   SubjectType = "10603"
	SubjectType10604   SubjectType = "10604"
	SubjectType10605   SubjectType = "10605"
	SubjectType107     SubjectType = "107"
	SubjectType10701   SubjectType = "10701"
	SubjectType10702   SubjectType = "10702"
	SubjectType108     SubjectType = "108"
	SubjectType10801   SubjectType = "10801"
	SubjectType10802   SubjectType = "10802"
	SubjectType10803   SubjectType = "10803"
	SubjectType12      SubjectType = "12"
	SubjectType109     SubjectType = "109"
	SubjectType10901   SubjectType = "10901"
	SubjectType10902   SubjectType = "10902"
	SubjectType10903   SubjectType = "10903"
	SubjectType110     SubjectType = "110"
	SubjectType11001   SubjectType = "11001"
	SubjectType11002   SubjectType = "11002"
	SubjectType11003   SubjectType = "11003"
	SubjectType11004   SubjectType = "11004"
	SubjectType111     SubjectType = "111"
	SubjectType11101   SubjectType = "11101"
	SubjectType11102   SubjectType = "11102"
	SubjectType11103   SubjectType = "11103"
	SubjectType11104   SubjectType = "11104"
	SubjectType112     SubjectType = "112"
	SubjectType11201   SubjectType = "11201"
	SubjectType11202   SubjectType = "11202"
	SubjectType11203   SubjectType = "11203"
	SubjectType11204   SubjectType = "11204"
	SubjectType11205   SubjectType = "11205"
	SubjectType11206y  SubjectType = "11206y"
	SubjectType113     SubjectType = "113"
	SubjectType11301   SubjectType = "11301"
	SubjectType11302   SubjectType = "11302"
	SubjectType11303   SubjectType = "11303"
	SubjectType11304   SubjectType = "11304"
	SubjectType11305   SubjectType = "11305"
	SubjectType2       SubjectType = "2"
	SubjectType21      SubjectType = "21"
	SubjectType201     SubjectType = "201"
	SubjectType20101   SubjectType = "20101"
	SubjectType20102   SubjectType = "20102"
	SubjectType20103   SubjectType = "20103"
	SubjectType20104   SubjectType = "20104"
	SubjectType20105   SubjectType = "20105"
	SubjectType20106   SubjectType = "20106"
	SubjectType20107   SubjectType = "20107"
	SubjectType20108   SubjectType = "20108"
	SubjectType202     SubjectType = "202"
	SubjectType20201   SubjectType = "20201"
	SubjectType20202   SubjectType = "20202"
	SubjectType20203   SubjectType = "20203"
	SubjectType20204   SubjectType = "20204"
	SubjectType20205   SubjectType = "20205"
	SubjectType20206   SubjectType = "20206"
	SubjectType20207   SubjectType = "20207"
	SubjectType203     SubjectType = "203"
	SubjectType20301   SubjectType = "20301"
	SubjectType20302   SubjectType = "20302"
	SubjectType20303   SubjectType = "20303"
	SubjectType20304   SubjectType = "20304"
	SubjectType20305   SubjectType = "20305"
	SubjectType20306   SubjectType = "20306"
	SubjectType22      SubjectType = "22"
	SubjectType204     SubjectType = "204"
	SubjectType20401   SubjectType = "20401"
	SubjectType20402   SubjectType = "20402"
	SubjectType20403   SubjectType = "20403"
	SubjectType20404   SubjectType = "20404"
	SubjectType20405   SubjectType = "20405"
	SubjectType205     SubjectType = "205"
	SubjectType20501cs SubjectType = "20501cs"
	SubjectType20502   SubjectType = "20502"
	SubjectType20503   SubjectType = "20503"
	SubjectType20504   SubjectType = "20504"
	SubjectType20505   SubjectType = "20505"
	SubjectType20506   SubjectType = "20506"
	SubjectType20507   SubjectType = "20507"
	SubjectType20508   SubjectType = "20508"
	SubjectType20509   SubjectType = "20509"
	SubjectType20510   SubjectType = "20510"
	SubjectType20511   SubjectType = "20511"
	SubjectType20512   SubjectType = "20512"
	SubjectType20513   SubjectType = "20513"
	SubjectType20514   SubjectType = "20514"
	SubjectType20515   SubjectType = "20515"
	SubjectType20516   SubjectType = "20516"
	SubjectType20517   SubjectType = "20517"
	SubjectType20518   SubjectType = "20518"
	SubjectType20519   SubjectType = "20519"
	SubjectType20520   SubjectType = "20520"
	SubjectType20521   SubjectType = "20521"
	SubjectType20522   SubjectType = "20522"
	SubjectType20523   SubjectType = "20523"
	SubjectType20524   SubjectType = "20524"
	SubjectType20525   SubjectType = "20525"
	SubjectType20526   SubjectType = "20526"
	SubjectType20527   SubjectType = "20527"
	SubjectType20528   SubjectType = "20528"
	SubjectType20529   SubjectType = "20529"
	SubjectType20530   SubjectType = "20530"
	SubjectType20531   SubjectType = "20531"
	SubjectType20532   SubjectType = "20532"
	SubjectType206     SubjectType = "206"
	SubjectType20601   SubjectType = "20601"
	SubjectType20602   SubjectType = "20602"
	SubjectType20603   SubjectType = "20603"
	SubjectType20604   SubjectType = "20604"
	SubjectType20605   SubjectType = "20605"
	SubjectType20606   SubjectType = "20606"
	SubjectType20607   SubjectType = "20607"
	SubjectType20608   SubjectType = "20608"
	SubjectType20609   SubjectType = "20609"
	SubjectType20610   SubjectType = "20610"
	SubjectType20611   SubjectType = "20611"
	SubjectType23      SubjectType = "23"
	SubjectType207     SubjectType = "207"
	SubjectType20701   SubjectType = "20701"
	SubjectType20702   SubjectType = "20702"
	SubjectType20703   SubjectType = "20703"
	SubjectType20704   SubjectType = "20704"
	SubjectType20705   SubjectType = "20705"
	SubjectType20706   SubjectType = "20706"
	SubjectType20707   SubjectType = "20707"
	SubjectType20708   SubjectType = "20708"
	SubjectType20709   SubjectType = "20709"
	SubjectType20710   SubjectType = "20710"
	SubjectType20711   SubjectType = "20711"
	SubjectType20712   SubjectType = "20712"
	SubjectType20713   SubjectType = "20713"
	SubjectType20714   SubjectType = "20714"
	SubjectType3       SubjectType = "3"
	SubjectType31      SubjectType = "31"
	SubjectType301     SubjectType = "301"
	SubjectType30101   SubjectType = "30101"
	SubjectType30102   SubjectType = "30102"
	SubjectType302     SubjectType = "302"
	SubjectType30201   SubjectType = "30201"
	SubjectType30202   SubjectType = "30202"
	SubjectType30203   SubjectType = "30203"
	SubjectType303     SubjectType = "303"
	SubjectType30301   SubjectType = "30301"
	SubjectType30302   SubjectType = "30302"
	SubjectType304     SubjectType = "304"
	SubjectType30401   SubjectType = "30401"
	SubjectType305     SubjectType = "305"
	SubjectType30501   SubjectType = "30501"
	SubjectType30502   SubjectType = "30502"
	SubjectType306     SubjectType = "306"
	SubjectType30601   SubjectType = "30601"
	SubjectType30602   SubjectType = "30602"
	SubjectType30603   SubjectType = "30603"
	SubjectType32      SubjectType = "32"
	SubjectType307     SubjectType = "307"
	SubjectType30701   SubjectType = "30701"
	SubjectType30702   SubjectType = "30702"
	SubjectType308     SubjectType = "308"
	SubjectType30801   SubjectType = "30801"
	SubjectType309     SubjectType = "309"
	SubjectType30901   SubjectType = "30901"
	SubjectType310     SubjectType = "310"
	SubjectType31001   SubjectType = "31001"
	SubjectType311     SubjectType = "311"
	SubjectType31101   SubjectType = "31101"
	SubjectType33      SubjectType = "33"
	SubjectType312     SubjectType = "312"
	SubjectType31201   SubjectType = "31201"
	SubjectType34      SubjectType = "34"
	SubjectType313     SubjectType = "313"
	SubjectType31301   SubjectType = "31301"
	SubjectType31302   SubjectType = "31302"
	SubjectType314     SubjectType = "314"
	SubjectType31401   SubjectType = "31401"
	SubjectType315     SubjectType = "315"
	SubjectType31501   SubjectType = "31501"
	SubjectType31502   SubjectType = "31502"
	SubjectType316     SubjectType = "316"
	SubjectType31601   SubjectType = "31601"
	SubjectType317     SubjectType = "317"
	SubjectType31701   SubjectType = "31701"
	SubjectType31702   SubjectType = "31702"
	SubjectType318     SubjectType = "318"
	SubjectType31801   SubjectType = "31801"
	SubjectType4       SubjectType = "4"
	SubjectType41      SubjectType = "41"
	SubjectType401     SubjectType = "401"
	SubjectType40101   SubjectType = "40101"
	SubjectType40102   SubjectType = "40102"
	SubjectType40103   SubjectType = "40103"
	SubjectType40104   SubjectType = "40104"
	SubjectType40105   SubjectType = "40105"
	SubjectType402     SubjectType = "402"
	SubjectType40201   SubjectType = "40201"
	SubjectType40202   SubjectType = "40202"
	SubjectType40203   SubjectType = "40203"
	SubjectType40204   SubjectType = "40204"
	SubjectType42      SubjectType = "42"
	SubjectType403     SubjectType = "403"
	SubjectType40301   SubjectType = "40301"
	SubjectType40302   SubjectType = "40302"
	SubjectType40303   SubjectType = "40303"
	SubjectType40304   SubjectType = "40304"
	SubjectType404     SubjectType = "404"
	SubjectType40401   SubjectType = "40401"
	SubjectType40402   SubjectType = "40402"
	SubjectType40403   SubjectType = "40403"
	SubjectType40404   SubjectType = "40404"
	SubjectType43      SubjectType = "43"
	SubjectType405     SubjectType = "405"
	SubjectType40501   SubjectType = "40501"
	SubjectType40502   SubjectType = "40502"
	SubjectType40503   SubjectType = "40503"
	SubjectType40504   SubjectType = "40504"
	SubjectType40505   SubjectType = "40505"
	SubjectType406     SubjectType = "406"
	SubjectType40601   SubjectType = "40601"
	SubjectType40602   SubjectType = "40602"
	SubjectType40603   SubjectType = "40603"
	SubjectType40604   SubjectType = "40604"
	SubjectType40605   SubjectType = "40605"
	SubjectType44      SubjectType = "44"
	SubjectType407     SubjectType = "407"
	SubjectType40701   SubjectType = "40701"
	SubjectType40702   SubjectType = "40702"
	SubjectType40703   SubjectType = "40703"
	SubjectType40704   SubjectType = "40704"
	SubjectType40705   SubjectType = "40705"
	SubjectType408     SubjectType = "408"
	SubjectType40801   SubjectType = "40801"
	SubjectType40802   SubjectType = "40802"
	SubjectType40803   SubjectType = "40803"
	SubjectType409     SubjectType = "409"
	SubjectType40901   SubjectType = "40901"
	SubjectType40902   SubjectType = "40902"
	SubjectType40903   SubjectType = "40903"
	SubjectType40904   SubjectType = "40904"
	SubjectType40905   SubjectType = "40905"
	SubjectType45      SubjectType = "45"
	SubjectType410     SubjectType = "410"
	SubjectType41001   SubjectType = "41001"
	SubjectType41002   SubjectType = "41002"
	SubjectType41003   SubjectType = "41003"
	SubjectType41004   SubjectType = "41004"
	SubjectType41005   SubjectType = "41005"
	SubjectType41006   SubjectType = "41006"
)

var SubjectTypeReverse = map[string]SubjectType{
	string(SubjectType1):       SubjectType1,
	string(SubjectType11):      SubjectType11,
	string(SubjectType101):     SubjectType101,
	string(SubjectType10101):   SubjectType10101,
	string(SubjectType10102):   SubjectType10102,
	string(SubjectType10103):   SubjectType10103,
	string(SubjectType10104):   SubjectType10104,
	string(SubjectType10105):   SubjectType10105,
	string(SubjectType102):     SubjectType102,
	string(SubjectType10201):   SubjectType10201,
	string(SubjectType10202):   SubjectType10202,
	string(SubjectType10203):   SubjectType10203,
	string(SubjectType10204):   SubjectType10204,
	string(SubjectType103):     SubjectType103,
	string(SubjectType10301):   SubjectType10301,
	string(SubjectType10302):   SubjectType10302,
	string(SubjectType10303):   SubjectType10303,
	string(SubjectType104):     SubjectType104,
	string(SubjectType10401):   SubjectType10401,
	string(SubjectType10402):   SubjectType10402,
	string(SubjectType10403):   SubjectType10403,
	string(SubjectType105):     SubjectType105,
	string(SubjectType10501):   SubjectType10501,
	string(SubjectType10502):   SubjectType10502,
	string(SubjectType10503):   SubjectType10503,
	string(SubjectType10504):   SubjectType10504,
	string(SubjectType106):     SubjectType106,
	string(SubjectType10601):   SubjectType10601,
	string(SubjectType10602):   SubjectType10602,
	string(SubjectType10603):   SubjectType10603,
	string(SubjectType10604):   SubjectType10604,
	string(SubjectType10605):   SubjectType10605,
	string(SubjectType107):     SubjectType107,
	string(SubjectType10701):   SubjectType10701,
	string(SubjectType10702):   SubjectType10702,
	string(SubjectType108):     SubjectType108,
	string(SubjectType10801):   SubjectType10801,
	string(SubjectType10802):   SubjectType10802,
	string(SubjectType10803):   SubjectType10803,
	string(SubjectType12):      SubjectType12,
	string(SubjectType109):     SubjectType109,
	string(SubjectType10901):   SubjectType10901,
	string(SubjectType10902):   SubjectType10902,
	string(SubjectType10903):   SubjectType10903,
	string(SubjectType110):     SubjectType110,
	string(SubjectType11001):   SubjectType11001,
	string(SubjectType11002):   SubjectType11002,
	string(SubjectType11003):   SubjectType11003,
	string(SubjectType11004):   SubjectType11004,
	string(SubjectType111):     SubjectType111,
	string(SubjectType11101):   SubjectType11101,
	string(SubjectType11102):   SubjectType11102,
	string(SubjectType11103):   SubjectType11103,
	string(SubjectType11104):   SubjectType11104,
	string(SubjectType112):     SubjectType112,
	string(SubjectType11201):   SubjectType11201,
	string(SubjectType11202):   SubjectType11202,
	string(SubjectType11203):   SubjectType11203,
	string(SubjectType11204):   SubjectType11204,
	string(SubjectType11205):   SubjectType11205,
	string(SubjectType11206y):  SubjectType11206y,
	string(SubjectType113):     SubjectType113,
	string(SubjectType11301):   SubjectType11301,
	string(SubjectType11302):   SubjectType11302,
	string(SubjectType11303):   SubjectType11303,
	string(SubjectType11304):   SubjectType11304,
	string(SubjectType11305):   SubjectType11305,
	string(SubjectType2):       SubjectType2,
	string(SubjectType21):      SubjectType21,
	string(SubjectType201):     SubjectType201,
	string(SubjectType20101):   SubjectType20101,
	string(SubjectType20102):   SubjectType20102,
	string(SubjectType20103):   SubjectType20103,
	string(SubjectType20104):   SubjectType20104,
	string(SubjectType20105):   SubjectType20105,
	string(SubjectType20106):   SubjectType20106,
	string(SubjectType20107):   SubjectType20107,
	string(SubjectType20108):   SubjectType20108,
	string(SubjectType202):     SubjectType202,
	string(SubjectType20201):   SubjectType20201,
	string(SubjectType20202):   SubjectType20202,
	string(SubjectType20203):   SubjectType20203,
	string(SubjectType20204):   SubjectType20204,
	string(SubjectType20205):   SubjectType20205,
	string(SubjectType20206):   SubjectType20206,
	string(SubjectType20207):   SubjectType20207,
	string(SubjectType203):     SubjectType203,
	string(SubjectType20301):   SubjectType20301,
	string(SubjectType20302):   SubjectType20302,
	string(SubjectType20303):   SubjectType20303,
	string(SubjectType20304):   SubjectType20304,
	string(SubjectType20305):   SubjectType20305,
	string(SubjectType20306):   SubjectType20306,
	string(SubjectType22):      SubjectType22,
	string(SubjectType204):     SubjectType204,
	string(SubjectType20401):   SubjectType20401,
	string(SubjectType20402):   SubjectType20402,
	string(SubjectType20403):   SubjectType20403,
	string(SubjectType20404):   SubjectType20404,
	string(SubjectType20405):   SubjectType20405,
	string(SubjectType205):     SubjectType205,
	string(SubjectType20501cs): SubjectType20501cs,
	string(SubjectType20502):   SubjectType20502,
	string(SubjectType20503):   SubjectType20503,
	string(SubjectType20504):   SubjectType20504,
	string(SubjectType20505):   SubjectType20505,
	string(SubjectType20506):   SubjectType20506,
	string(SubjectType20507):   SubjectType20507,
	string(SubjectType20508):   SubjectType20508,
	string(SubjectType20509):   SubjectType20509,
	string(SubjectType20510):   SubjectType20510,
	string(SubjectType20511):   SubjectType20511,
	string(SubjectType20512):   SubjectType20512,
	string(SubjectType20513):   SubjectType20513,
	string(SubjectType20514):   SubjectType20514,
	string(SubjectType20515):   SubjectType20515,
	string(SubjectType20516):   SubjectType20516,
	string(SubjectType20517):   SubjectType20517,
	string(SubjectType20518):   SubjectType20518,
	string(SubjectType20519):   SubjectType20519,
	string(SubjectType20520):   SubjectType20520,
	string(SubjectType20521):   SubjectType20521,
	string(SubjectType20522):   SubjectType20522,
	string(SubjectType20523):   SubjectType20523,
	string(SubjectType20524):   SubjectType20524,
	string(SubjectType20525):   SubjectType20525,
	string(SubjectType20526):   SubjectType20526,
	string(SubjectType20527):   SubjectType20527,
	string(SubjectType20528):   SubjectType20528,
	string(SubjectType20529):   SubjectType20529,
	string(SubjectType20530):   SubjectType20530,
	string(SubjectType20531):   SubjectType20531,
	string(SubjectType20532):   SubjectType20532,
	string(SubjectType206):     SubjectType206,
	string(SubjectType20601):   SubjectType20601,
	string(SubjectType20602):   SubjectType20602,
	string(SubjectType20603):   SubjectType20603,
	string(SubjectType20604):   SubjectType20604,
	string(SubjectType20605):   SubjectType20605,
	string(SubjectType20606):   SubjectType20606,
	string(SubjectType20607):   SubjectType20607,
	string(SubjectType20608):   SubjectType20608,
	string(SubjectType20609):   SubjectType20609,
	string(SubjectType20610):   SubjectType20610,
	string(SubjectType20611):   SubjectType20611,
	string(SubjectType23):      SubjectType23,
	string(SubjectType207):     SubjectType207,
	string(SubjectType20701):   SubjectType20701,
	string(SubjectType20702):   SubjectType20702,
	string(SubjectType20703):   SubjectType20703,
	string(SubjectType20704):   SubjectType20704,
	string(SubjectType20705):   SubjectType20705,
	string(SubjectType20706):   SubjectType20706,
	string(SubjectType20707):   SubjectType20707,
	string(SubjectType20708):   SubjectType20708,
	string(SubjectType20709):   SubjectType20709,
	string(SubjectType20710):   SubjectType20710,
	string(SubjectType20711):   SubjectType20711,
	string(SubjectType20712):   SubjectType20712,
	string(SubjectType20713):   SubjectType20713,
	string(SubjectType20714):   SubjectType20714,
	string(SubjectType3):       SubjectType3,
	string(SubjectType31):      SubjectType31,
	string(SubjectType301):     SubjectType301,
	string(SubjectType30101):   SubjectType30101,
	string(SubjectType30102):   SubjectType30102,
	string(SubjectType302):     SubjectType302,
	string(SubjectType30201):   SubjectType30201,
	string(SubjectType30202):   SubjectType30202,
	string(SubjectType30203):   SubjectType30203,
	string(SubjectType303):     SubjectType303,
	string(SubjectType30301):   SubjectType30301,
	string(SubjectType30302):   SubjectType30302,
	string(SubjectType304):     SubjectType304,
	string(SubjectType30401):   SubjectType30401,
	string(SubjectType305):     SubjectType305,
	string(SubjectType30501):   SubjectType30501,
	string(SubjectType30502):   SubjectType30502,
	string(SubjectType306):     SubjectType306,
	string(SubjectType30601):   SubjectType30601,
	string(SubjectType30602):   SubjectType30602,
	string(SubjectType30603):   SubjectType30603,
	string(SubjectType32):      SubjectType32,
	string(SubjectType307):     SubjectType307,
	string(SubjectType30701):   SubjectType30701,
	string(SubjectType30702):   SubjectType30702,
	string(SubjectType308):     SubjectType308,
	string(SubjectType30801):   SubjectType30801,
	string(SubjectType309):     SubjectType309,
	string(SubjectType30901):   SubjectType30901,
	string(SubjectType310):     SubjectType310,
	string(SubjectType31001):   SubjectType31001,
	string(SubjectType311):     SubjectType311,
	string(SubjectType31101):   SubjectType31101,
	string(SubjectType33):      SubjectType33,
	string(SubjectType312):     SubjectType312,
	string(SubjectType31201):   SubjectType31201,
	string(SubjectType34):      SubjectType34,
	string(SubjectType313):     SubjectType313,
	string(SubjectType31301):   SubjectType31301,
	string(SubjectType31302):   SubjectType31302,
	string(SubjectType314):     SubjectType314,
	string(SubjectType31401):   SubjectType31401,
	string(SubjectType315):     SubjectType315,
	string(SubjectType31501):   SubjectType31501,
	string(SubjectType31502):   SubjectType31502,
	string(SubjectType316):     SubjectType316,
	string(SubjectType31601):   SubjectType31601,
	string(SubjectType317):     SubjectType317,
	string(SubjectType31701):   SubjectType31701,
	string(SubjectType31702):   SubjectType31702,
	string(SubjectType318):     SubjectType318,
	string(SubjectType31801):   SubjectType31801,
	string(SubjectType4):       SubjectType4,
	string(SubjectType41):      SubjectType41,
	string(SubjectType401):     SubjectType401,
	string(SubjectType40101):   SubjectType40101,
	string(SubjectType40102):   SubjectType40102,
	string(SubjectType40103):   SubjectType40103,
	string(SubjectType40104):   SubjectType40104,
	string(SubjectType40105):   SubjectType40105,
	string(SubjectType402):     SubjectType402,
	string(SubjectType40201):   SubjectType40201,
	string(SubjectType40202):   SubjectType40202,
	string(SubjectType40203):   SubjectType40203,
	string(SubjectType40204):   SubjectType40204,
	string(SubjectType42):      SubjectType42,
	string(SubjectType403):     SubjectType403,
	string(SubjectType40301):   SubjectType40301,
	string(SubjectType40302):   SubjectType40302,
	string(SubjectType40303):   SubjectType40303,
	string(SubjectType40304):   SubjectType40304,
	string(SubjectType404):     SubjectType404,
	string(SubjectType40401):   SubjectType40401,
	string(SubjectType40402):   SubjectType40402,
	string(SubjectType40403):   SubjectType40403,
	string(SubjectType40404):   SubjectType40404,
	string(SubjectType43):      SubjectType43,
	string(SubjectType405):     SubjectType405,
	string(SubjectType40501):   SubjectType40501,
	string(SubjectType40502):   SubjectType40502,
	string(SubjectType40503):   SubjectType40503,
	string(SubjectType40504):   SubjectType40504,
	string(SubjectType40505):   SubjectType40505,
	string(SubjectType406):     SubjectType406,
	string(SubjectType40601):   SubjectType40601,
	string(SubjectType40602):   SubjectType40602,
	string(SubjectType40603):   SubjectType40603,
	string(SubjectType40604):   SubjectType40604,
	string(SubjectType40605):   SubjectType40605,
	string(SubjectType44):      SubjectType44,
	string(SubjectType407):     SubjectType407,
	string(SubjectType40701):   SubjectType40701,
	string(SubjectType40702):   SubjectType40702,
	string(SubjectType40703):   SubjectType40703,
	string(SubjectType40704):   SubjectType40704,
	string(SubjectType40705):   SubjectType40705,
	string(SubjectType408):     SubjectType408,
	string(SubjectType40801):   SubjectType40801,
	string(SubjectType40802):   SubjectType40802,
	string(SubjectType40803):   SubjectType40803,
	string(SubjectType409):     SubjectType409,
	string(SubjectType40901):   SubjectType40901,
	string(SubjectType40902):   SubjectType40902,
	string(SubjectType40903):   SubjectType40903,
	string(SubjectType40904):   SubjectType40904,
	string(SubjectType40905):   SubjectType40905,
	string(SubjectType45):      SubjectType45,
	string(SubjectType410):     SubjectType410,
	string(SubjectType41001):   SubjectType41001,
	string(SubjectType41002):   SubjectType41002,
	string(SubjectType41003):   SubjectType41003,
	string(SubjectType41004):   SubjectType41004,
	string(SubjectType41005):   SubjectType41005,
	string(SubjectType41006):   SubjectType41006,
}

var SubjectTypeName = map[SubjectType]string{
	SubjectType1:       "Humanities and Social Sciences",
	SubjectType11:      "Humanities",
	SubjectType101:     "Ancient Cultures",
	SubjectType10101:   "Prehistory",
	SubjectType10102:   "Classical Philology",
	SubjectType10103:   "Ancient History",
	SubjectType10104:   "Classical Archaeology",
	SubjectType10105:   "Egyptology and Ancient Near Eastern Studies",
	SubjectType102:     "History",
	SubjectType10201:   "Medieval History",
	SubjectType10202:   "Early Modern History",
	SubjectType10203:   "Modern and Current History",
	SubjectType10204:   "History of Science",
	SubjectType103:     "Fine Arts, Music, Theatre and Media Studies",
	SubjectType10301:   "Art History",
	SubjectType10302:   "Musicology",
	SubjectType10303:   "Theatre and Media Studies",
	SubjectType104:     "Linguistics",
	SubjectType10401:   "General and Applied Linguistics",
	SubjectType10402:   "Individual Linguistics",
	SubjectType10403:   "Typology, Non-European Languages, Historical Linguistics",
	SubjectType105:     "Literary Studies",
	SubjectType10501:   "Medieval German Literature",
	SubjectType10502:   "Modern German Literature",
	SubjectType10503:   "European and American Literature",
	SubjectType10504:   "General and Comparative Literature and Cultural Studies",
	SubjectType106:     "Non-European Languages and Cultures, Social and Cultural Anthropology, Jewish Studies and Religious Studies",
	SubjectType10601:   "Social and Cultural Anthropology and Ethnology/Folklore",
	SubjectType10602:   "Asian Studies",
	SubjectType10603:   "African, American and Oceania Studies",
	SubjectType10604:   "Islamic Studies, Arabian Studies, Semitic Studies",
	SubjectType10605:   "Religious Studies and Jewish Studies",
	SubjectType107:     "Theology",
	SubjectType10701:   "Protestant Theology",
	SubjectType10702:   "Roman Catholic Theology",
	SubjectType108:     "Philosophy",
	SubjectType10801:   "History of Philosophy",
	SubjectType10802:   "Theoretical Philosophy",
	SubjectType10803:   "Practical Philosophy",
	SubjectType12:      "Social and Behavioural Sciences",
	SubjectType109:     "Education Sciences",
	SubjectType10901:   "General Education and History of Education",
	SubjectType10902:   "Research on Teaching, Learning and Training",
	SubjectType10903:   "Research on Socialization and Educational Institutions and Professions",
	SubjectType110:     "Psychology",
	SubjectType11001:   "General, Biological and Mathematical Psychology",
	SubjectType11002:   "Developmental and Educational Psychology",
	SubjectType11003:   "Social Psychology, Industrial and Organisational Psychology",
	SubjectType11004:   "Differential Psychology, Clinical Psychology, Medical Psychology, Methodology",
	SubjectType111:     "Social Sciences",
	SubjectType11101:   "Sociological Theory",
	SubjectType11102:   "Empirical Social Research",
	SubjectType11103:   "Communication Science",
	SubjectType11104:   "Political Science",
	SubjectType112:     "Economics",
	SubjectType11201:   "Economic Theory",
	SubjectType11202:   "Economic and Social Policy",
	SubjectType11203:   "Public Finance",
	SubjectType11204:   "Business Administration",
	SubjectType11205:   "Statistics and Econometrics",
	SubjectType11206y:  "Economic and Social History",
	SubjectType113:     "Jurisprudence",
	SubjectType11301:   "Legal and Political Philosophy, Legal History, Legal Theory",
	SubjectType11302:   "Private Law",
	SubjectType11303:   "Public Law",
	SubjectType11304:   "Criminal Law and Law of Criminal Procedure",
	SubjectType11305:   "Criminology",
	SubjectType2:       "Life Sciences",
	SubjectType21:      "Biology",
	SubjectType201:     "Basic Biological and Medical Research",
	SubjectType20101:   "Biochemistry",
	SubjectType20102:   "Biophysics",
	SubjectType20103:   "Cell Biology",
	SubjectType20104:   "Structural Biology",
	SubjectType20105:   "General Genetics",
	SubjectType20106:   "Developmental Biology",
	SubjectType20107:   "Bioinformatics and Theoretical Biology",
	SubjectType20108:   "Anatomy",
	SubjectType202:     "Plant Sciences",
	SubjectType20201:   "Plant Systematics and Evolution",
	SubjectType20202:   "Plant Ecology and Ecosystem Analysis",
	SubjectType20203:   "Inter-organismic Interactions of Plants",
	SubjectType20204:   "Plant Physiology",
	SubjectType20205:   "Plant Biochemistry and Biophysics",
	SubjectType20206:   "Plant Cell and Developmental Biology",
	SubjectType20207:   "Plant Genetics",
	SubjectType203:     "Zoology",
	SubjectType20301:   "Systematics and Morphology",
	SubjectType20302:   "Evolution, Anthropology",
	SubjectType20303:   "Animal Ecology, Biodiversity and Ecosystem Research",
	SubjectType20304:   "Sensory and Behavioural Biology",
	SubjectType20305:   "Biochemistry and Animal Physiology",
	SubjectType20306:   "Animal Genetics, Cell and Developmental Biology",
	SubjectType22:      "Medicine",
	SubjectType204:     "Microbiology, Virology and Immunology",
	SubjectType20401:   "Metabolism, Biochemistry and Genetics of Microorganisms",
	SubjectType20402:   "Microbial Ecology and Applied Microbiology",
	SubjectType20403:   "Medical Microbiology, Molecular Infection Biology",
	SubjectType20404:   "Virology",
	SubjectType20405:   "Immunology",
	SubjectType205:     "Medicine",
	SubjectType20501cs: "Epidemiology, Medical Biometry, Medical Informatics",
	SubjectType20502:   "Public Health, Health Services Research, Social Medicine",
	SubjectType20503:   "Human Genetics",
	SubjectType20504:   "Physiology",
	SubjectType20505:   "Nutritional Sciences",
	SubjectType20506:   "Pathology and Forensic Medicine",
	SubjectType20507:   "Clinical Chemistry and Pathobiochemistry",
	SubjectType20508:   "Pharmacy",
	SubjectType20509:   "Pharmacology",
	SubjectType20510:   "Toxicology and Occupational Medicine",
	SubjectType20511:   "Anaesthesiology",
	SubjectType20512:   "Cardiology, Angiology",
	SubjectType20513:   "Pneumology, Clinical Infectiology Intensive Care Medicine",
	SubjectType20514:   "Hematology, Oncology, Transfusion Medicine",
	SubjectType20515:   "Gastroenterology, Metabolism",
	SubjectType20516:   "Nephrology",
	SubjectType20517:   "Endocrinology, Diabetology",
	SubjectType20518:   "Rheumatology, Clinical Immunology, Allergology",
	SubjectType20519:   "Dermatology",
	SubjectType20520:   "Pediatric and Adolescent Medicine",
	SubjectType20521:   "Gynaecology and Obstetrics",
	SubjectType20522:   "Reproductive Medicine/Biology",
	SubjectType20523:   "Urology",
	SubjectType20524:   "Gerontology and Geriatric Medicine",
	SubjectType20525:   "Vascular and Visceral Surgery",
	SubjectType20526:   "Cardiothoracic Surgery",
	SubjectType20527:   "Traumatology and Orthopaedics",
	SubjectType20528:   "Dentistry, Oral Surgery",
	SubjectType20529:   "Otolaryngology",
	SubjectType20530:   "Radiology and Nuclear Medicine",
	SubjectType20531:   "Radiation Oncology and Radiobiology",
	SubjectType20532:   "Biomedical Technology and Medical Physics",
	SubjectType206:     "Neurosciences",
	SubjectType20601:   "Molecular Neuroscience and Neurogenetics",
	SubjectType20602:   "Cellular Neuroscience",
	SubjectType20603:   "Developmental Neurobiology",
	SubjectType20604:   "Systemic Neuroscience, Computational Neuroscience, Behaviour",
	SubjectType20605:   "Comparative Neurobiology",
	SubjectType20606:   "Cognitive Neuroscience and Neuroimaging",
	SubjectType20607:   "Molecular Neurology",
	SubjectType20608:   "Clinical Neurosciences I - Neurology, Neurosurgery",
	SubjectType20609:   "Biological Psychiatry",
	SubjectType20610:   "Clinical Neurosciences II - Psychiatry, Psychotherapy, Psychosomatic Medicine",
	SubjectType20611:   "Clinical Neurosciences III - Ophthalmology",
	SubjectType23:      "Agriculture, Forestry, Horticulture and Veterinary Medicine",
	SubjectType207:     "Agriculture, Forestry, Horticulture and Veterinary Medicine",
	SubjectType20701:   "Soil Sciences",
	SubjectType20702:   "Plant Cultivation",
	SubjectType20703:   "Plant Nutrition",
	SubjectType20704:   "Ecology of Agricultural Landscapes",
	SubjectType20705:   "Plant Breeding",
	SubjectType20706:   "Phytomedicine",
	SubjectType20707:   "Agricultural and Food Process Engineering",
	SubjectType20708:   "Agricultural Economics and Sociology",
	SubjectType20709:   "Inventory Control and Use of Forest Resources",
	SubjectType20710:   "Basic Forest Research",
	SubjectType20711:   "Animal Husbandry, Breeding and Hygiene",
	SubjectType20712:   "Animal Nutrition and Nutrition Physiology",
	SubjectType20713:   "Basic Veterinary Medical Science",
	SubjectType20714:   "Basic Research on Pathogenesis, Diagnostics and Therapy and Clinical Veterinary Medicine",
	SubjectType3:       "Natural Sciences",
	SubjectType31:      "Chemistry",
	SubjectType301:     "Molecular Chemistry",
	SubjectType30101:   "Inorganic Molecular Chemistry",
	SubjectType30102:   "Organic Molecular Chemistry",
	SubjectType302:     "Chemical Solid State and Surface Research",
	SubjectType30201:   "Solid State and Surface Chemistry, Material Synthesis",
	SubjectType30202:   "Physical Chemistry of Solids and Surfaces, Material Characterisation",
	SubjectType30203:   "Theory and Modelling",
	SubjectType303:     "Physical and Theoretical Chemistry",
	SubjectType30301:   "Physical Chemistry of Molecules, Interfaces and Liquids - Spectroscopy, Kinetics",
	SubjectType30302:   "General Theoretical Chemistry",
	SubjectType304:     "Analytical Chemistry, Method Development (Chemistry)",
	SubjectType30401:   "Analytical Chemistry, Method Development (Chemistry)",
	SubjectType305:     "Biological Chemistry and Food Chemistry",
	SubjectType30501:   "Biological and Biomimetic Chemistry",
	SubjectType30502:   "Food Chemistry",
	SubjectType306:     "Polymer Research",
	SubjectType30601:   "Preparatory and Physical Chemistry of Polymers",
	SubjectType30602:   "Experimental and Theoretical Physics of Polymers",
	SubjectType30603:   "Polymer Materials",
	SubjectType32:      "Physics",
	SubjectType307:     "Condensed Matter Physics",
	SubjectType30701:   "Experimental Condensed Matter Physics",
	SubjectType30702:   "Theoretical Condensed Matter Physics",
	SubjectType308:     "Optics, Quantum Optics and Physics of Atoms, Molecules and Plasmas",
	SubjectType30801:   "Optics, Quantum Optics, Atoms, Molecules, Plasmas",
	SubjectType309:     "Particles, Nuclei and Fields",
	SubjectType30901:   "Particles, Nuclei and Fields",
	SubjectType310:     "Statistical Physics, Soft Matter, Biological Physics, Nonlinear Dynamics",
	SubjectType31001:   "Statistical Physics, Soft Matter, Biological Physics, Nonlinear Dynamics",
	SubjectType311:     "Astrophysics and Astronomy",
	SubjectType31101:   "Astrophysics and Astronomy",
	SubjectType33:      "Mathematics",
	SubjectType312:     "Mathematics",
	SubjectType31201:   "Mathematics",
	SubjectType34:      "Geosciences (including Geography)",
	SubjectType313:     "Atmospheric Science and Oceanography",
	SubjectType31301:   "Atmospheric Science",
	SubjectType31302:   "Oceanography",
	SubjectType314:     "Geology and Palaeontology",
	SubjectType31401:   "Geology and Palaeontology",
	SubjectType315:     "Geophysics and Geodesy",
	SubjectType31501:   "Geophysics",
	SubjectType31502:   "Geodesy, Photogrammetry, Remote Sensing, Geoinformatics, Cartogaphy",
	SubjectType316:     "Geochemistry, Mineralogy and Crystallography",
	SubjectType31601:   "Geochemistry, Mineralogy and Crystallography",
	SubjectType317:     "Geography",
	SubjectType31701:   "Physical Geography",
	SubjectType31702:   "Human Geography",
	SubjectType318:     "Water Research",
	SubjectType31801:   "Hydrogeology, Hydrology, Limnology, Urban Water Management, Water Chemistry, Integrated Water Resources Management",
	SubjectType4:       "Engineering Sciences",
	SubjectType41:      "Mechanical and industrial Engineering",
	SubjectType401:     "Production Technology",
	SubjectType40101:   "Metal-Cutting Manufacturing Engineering",
	SubjectType40102:   "Primary Shaping and Reshaping Technology",
	SubjectType40103:   "Micro-, Precision, Mounting, Joining, Separation Technology",
	SubjectType40104:   "Plastics Engineering",
	SubjectType40105:   "Production Automation, Factory Operation, Operations Manangement",
	SubjectType402:     "Mechanics and Constructive Mechanical Engineering",
	SubjectType40201:   "Construction, Machine Elements",
	SubjectType40202:   "Mechanics",
	SubjectType40203:   "Lightweight Construction, Textile Technology",
	SubjectType40204:   "Acoustics",
	SubjectType42:      "Thermal Engineering/Process Engineering",
	SubjectType403:     "Process Engineering, Technical Chemistry",
	SubjectType40301:   "Chemical and Thermal Process Engineering",
	SubjectType40302:   "Technical Chemistry",
	SubjectType40303:   "Mechanical Process Engineering",
	SubjectType40304:   "Biological Process Engineering",
	SubjectType404:     "Heat Energy Technology, Thermal Machines, Fluid Mechanics",
	SubjectType40401:   "Energy Process Engineering",
	SubjectType40402:   "Technical Thermodynamics",
	SubjectType40403:   "Fluid Mechanics",
	SubjectType40404:   "Hydraulic and Turbo Engines and Piston Engines",
	SubjectType43:      "Materials Science and Engineering",
	SubjectType405:     "Materials Engineering",
	SubjectType40501:   "Metallurgical and Thermal Processes, Thermomechanical Treatment of Materials",
	SubjectType40502:   "Sintered Metallic and Ceramic Materials",
	SubjectType40503:   "Composite Materials",
	SubjectType40504:   "Mechanical Behaviour of Construction Materials",
	SubjectType40505:   "Coating and Surface Technology",
	SubjectType406:     "Materials Science",
	SubjectType40601:   "Thermodynamics and Kinetics of Materials",
	SubjectType40602:   "Synthesis and Properties of Functional Materials",
	SubjectType40603:   "Microstructural Mechanical Properties of Materials",
	SubjectType40604:   "Structuring and Functionalisation",
	SubjectType40605:   "Biomaterials",
	SubjectType44:      "Computer Science, Electrical and System Engineering",
	SubjectType407:     "Systems Engineering",
	SubjectType40701:   "Automation, Control Systems, Robotics, Mechatronics",
	SubjectType40702:   "Measurement Systems",
	SubjectType40703:   "Microsystems",
	SubjectType40704:   "Traffic and Transport Systems, Logistics",
	SubjectType40705:   "Human Factors, Ergonomics, Human-Machine Systems",
	SubjectType408:     "Electrical Engineering",
	SubjectType40801:   "Electronic Semiconductors, Components, Circuits, Systems",
	SubjectType40802:   "Communication, High-Frequency and Network Technology, Theoretical Electrical Engineering",
	SubjectType40803:   "Electrical Energy Generation, Distribution, Application",
	SubjectType409:     "Computer Science",
	SubjectType40901:   "Theoretical Computer Science",
	SubjectType40902:   "Software Technology",
	SubjectType40903:   "Operating, Communication and Information Systems",
	SubjectType40904:   "Artificial Intelligence, Image and Language Processing",
	SubjectType40905:   "Computer Architecture and Embedded Systems",
	SubjectType45:      "Construction Engineering and Architecture",
	SubjectType410:     "Construction Engineering and Architecture",
	SubjectType41001:   "Architecture, Building and Construction History, Sustainable Building Technology, Building Design",
	SubjectType41002:   "Urbanism, Spatial Planning, Transportation and Infrastructure Planning, Landscape Planning",
	SubjectType41003:   "Construction Material Sciences, Chemistry, Building Physics",
	SubjectType41004:   "Structural Engineering, Building Informatics, Construction Operation",
	SubjectType41005:   "Applied Mechanics, Statics and Dynamics",
	SubjectType41006:   "Geotechnics, Hydraulic Engineering",
}

type ProviderType string

const (
	ProviderTypeDataProvider    ProviderType = "dataProvider"
	ProviderTypeServiceProvider ProviderType = "serviceProvider"
)

var ProviderTypeReverse = map[string]ProviderType{
	string(ProviderTypeDataProvider):    ProviderTypeDataProvider,
	string(ProviderTypeServiceProvider): ProviderTypeServiceProvider,
}

type AccessType string

const (
	AccessTypeOpen       AccessType = "open"
	AccessTypeEmbargoed  AccessType = "embargoed"
	AccessTypeRestricted AccessType = "restricted"
	AccessTypeClosed     AccessType = "closed"
)

var AccessTypeReverse = map[string]AccessType{
	string(AccessTypeOpen):       AccessTypeOpen,
	string(AccessTypeEmbargoed):  AccessTypeEmbargoed,
	string(AccessTypeRestricted): AccessTypeRestricted,
	string(AccessTypeClosed):     AccessTypeClosed,
}

type AccessRestrictions string

const (
	AccessRestrictionsFeeRequired             AccessRestrictions = "feeRequired"
	AccessRestrictionsInstitutionalMembership AccessRestrictions = "institutional membership"
	AccessRestrictionsRegistration            AccessRestrictions = "registration"
	AccessRestrictionsFeeOther                AccessRestrictions = "other"
)

var AccessRestrictionsReverse = map[string]AccessRestrictions{
	string(AccessRestrictionsFeeRequired):             AccessRestrictionsFeeRequired,
	string(AccessRestrictionsInstitutionalMembership): AccessRestrictionsInstitutionalMembership,
	string(AccessRestrictionsRegistration):            AccessRestrictionsRegistration,
	string(AccessRestrictionsFeeOther):                AccessRestrictionsFeeOther,
}

type YesNoUn string

const (
	YesNoUnYes         YesNoUn = "yes"
	DataYesNoUnNo      YesNoUn = "no"
	DataYesNoUnUnknown YesNoUn = "unknown"
)

var YesNoUnReverse = map[string]YesNoUn{
	string(YesNoUnYes):         YesNoUnYes,
	string(DataYesNoUnNo):      DataYesNoUnNo,
	string(DataYesNoUnUnknown): DataYesNoUnUnknown,
}
