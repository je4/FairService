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

type RE3DataSubjectType string

const (
	RE3DataSubjectType1       RE3DataSubjectType = "1"
	RE3DataSubjectType11      RE3DataSubjectType = "11"
	RE3DataSubjectType101     RE3DataSubjectType = "101"
	RE3DataSubjectType10101   RE3DataSubjectType = "10101"
	RE3DataSubjectType10102   RE3DataSubjectType = "10102"
	RE3DataSubjectType10103   RE3DataSubjectType = "10103"
	RE3DataSubjectType10104   RE3DataSubjectType = "10104"
	RE3DataSubjectType10105   RE3DataSubjectType = "10105"
	RE3DataSubjectType102     RE3DataSubjectType = "102"
	RE3DataSubjectType10201   RE3DataSubjectType = "10201"
	RE3DataSubjectType10202   RE3DataSubjectType = "10202"
	RE3DataSubjectType10203   RE3DataSubjectType = "10203"
	RE3DataSubjectType10204   RE3DataSubjectType = "10204"
	RE3DataSubjectType103     RE3DataSubjectType = "103"
	RE3DataSubjectType10301   RE3DataSubjectType = "10301"
	RE3DataSubjectType10302   RE3DataSubjectType = "10302"
	RE3DataSubjectType10303   RE3DataSubjectType = "10303"
	RE3DataSubjectType104     RE3DataSubjectType = "104"
	RE3DataSubjectType10401   RE3DataSubjectType = "10401"
	RE3DataSubjectType10402   RE3DataSubjectType = "10402"
	RE3DataSubjectType10403   RE3DataSubjectType = "10403"
	RE3DataSubjectType105     RE3DataSubjectType = "105"
	RE3DataSubjectType10501   RE3DataSubjectType = "10501"
	RE3DataSubjectType10502   RE3DataSubjectType = "10502"
	RE3DataSubjectType10503   RE3DataSubjectType = "10503"
	RE3DataSubjectType10504   RE3DataSubjectType = "10504"
	RE3DataSubjectType106     RE3DataSubjectType = "106"
	RE3DataSubjectType10601   RE3DataSubjectType = "10601"
	RE3DataSubjectType10602   RE3DataSubjectType = "10602"
	RE3DataSubjectType10603   RE3DataSubjectType = "10603"
	RE3DataSubjectType10604   RE3DataSubjectType = "10604"
	RE3DataSubjectType10605   RE3DataSubjectType = "10605"
	RE3DataSubjectType107     RE3DataSubjectType = "107"
	RE3DataSubjectType10701   RE3DataSubjectType = "10701"
	RE3DataSubjectType10702   RE3DataSubjectType = "10702"
	RE3DataSubjectType108     RE3DataSubjectType = "108"
	RE3DataSubjectType10801   RE3DataSubjectType = "10801"
	RE3DataSubjectType10802   RE3DataSubjectType = "10802"
	RE3DataSubjectType10803   RE3DataSubjectType = "10803"
	RE3DataSubjectType12      RE3DataSubjectType = "12"
	RE3DataSubjectType109     RE3DataSubjectType = "109"
	RE3DataSubjectType10901   RE3DataSubjectType = "10901"
	RE3DataSubjectType10902   RE3DataSubjectType = "10902"
	RE3DataSubjectType10903   RE3DataSubjectType = "10903"
	RE3DataSubjectType110     RE3DataSubjectType = "110"
	RE3DataSubjectType11001   RE3DataSubjectType = "11001"
	RE3DataSubjectType11002   RE3DataSubjectType = "11002"
	RE3DataSubjectType11003   RE3DataSubjectType = "11003"
	RE3DataSubjectType11004   RE3DataSubjectType = "11004"
	RE3DataSubjectType111     RE3DataSubjectType = "111"
	RE3DataSubjectType11101   RE3DataSubjectType = "11101"
	RE3DataSubjectType11102   RE3DataSubjectType = "11102"
	RE3DataSubjectType11103   RE3DataSubjectType = "11103"
	RE3DataSubjectType11104   RE3DataSubjectType = "11104"
	RE3DataSubjectType112     RE3DataSubjectType = "112"
	RE3DataSubjectType11201   RE3DataSubjectType = "11201"
	RE3DataSubjectType11202   RE3DataSubjectType = "11202"
	RE3DataSubjectType11203   RE3DataSubjectType = "11203"
	RE3DataSubjectType11204   RE3DataSubjectType = "11204"
	RE3DataSubjectType11205   RE3DataSubjectType = "11205"
	RE3DataSubjectType11206y  RE3DataSubjectType = "11206y"
	RE3DataSubjectType113     RE3DataSubjectType = "113"
	RE3DataSubjectType11301   RE3DataSubjectType = "11301"
	RE3DataSubjectType11302   RE3DataSubjectType = "11302"
	RE3DataSubjectType11303   RE3DataSubjectType = "11303"
	RE3DataSubjectType11304   RE3DataSubjectType = "11304"
	RE3DataSubjectType11305   RE3DataSubjectType = "11305"
	RE3DataSubjectType2       RE3DataSubjectType = "2"
	RE3DataSubjectType21      RE3DataSubjectType = "21"
	RE3DataSubjectType201     RE3DataSubjectType = "201"
	RE3DataSubjectType20101   RE3DataSubjectType = "20101"
	RE3DataSubjectType20102   RE3DataSubjectType = "20102"
	RE3DataSubjectType20103   RE3DataSubjectType = "20103"
	RE3DataSubjectType20104   RE3DataSubjectType = "20104"
	RE3DataSubjectType20105   RE3DataSubjectType = "20105"
	RE3DataSubjectType20106   RE3DataSubjectType = "20106"
	RE3DataSubjectType20107   RE3DataSubjectType = "20107"
	RE3DataSubjectType20108   RE3DataSubjectType = "20108"
	RE3DataSubjectType202     RE3DataSubjectType = "202"
	RE3DataSubjectType20201   RE3DataSubjectType = "20201"
	RE3DataSubjectType20202   RE3DataSubjectType = "20202"
	RE3DataSubjectType20203   RE3DataSubjectType = "20203"
	RE3DataSubjectType20204   RE3DataSubjectType = "20204"
	RE3DataSubjectType20205   RE3DataSubjectType = "20205"
	RE3DataSubjectType20206   RE3DataSubjectType = "20206"
	RE3DataSubjectType20207   RE3DataSubjectType = "20207"
	RE3DataSubjectType203     RE3DataSubjectType = "203"
	RE3DataSubjectType20301   RE3DataSubjectType = "20301"
	RE3DataSubjectType20302   RE3DataSubjectType = "20302"
	RE3DataSubjectType20303   RE3DataSubjectType = "20303"
	RE3DataSubjectType20304   RE3DataSubjectType = "20304"
	RE3DataSubjectType20305   RE3DataSubjectType = "20305"
	RE3DataSubjectType20306   RE3DataSubjectType = "20306"
	RE3DataSubjectType22      RE3DataSubjectType = "22"
	RE3DataSubjectType204     RE3DataSubjectType = "204"
	RE3DataSubjectType20401   RE3DataSubjectType = "20401"
	RE3DataSubjectType20402   RE3DataSubjectType = "20402"
	RE3DataSubjectType20403   RE3DataSubjectType = "20403"
	RE3DataSubjectType20404   RE3DataSubjectType = "20404"
	RE3DataSubjectType20405   RE3DataSubjectType = "20405"
	RE3DataSubjectType205     RE3DataSubjectType = "205"
	RE3DataSubjectType20501cs RE3DataSubjectType = "20501cs"
	RE3DataSubjectType20502   RE3DataSubjectType = "20502"
	RE3DataSubjectType20503   RE3DataSubjectType = "20503"
	RE3DataSubjectType20504   RE3DataSubjectType = "20504"
	RE3DataSubjectType20505   RE3DataSubjectType = "20505"
	RE3DataSubjectType20506   RE3DataSubjectType = "20506"
	RE3DataSubjectType20507   RE3DataSubjectType = "20507"
	RE3DataSubjectType20508   RE3DataSubjectType = "20508"
	RE3DataSubjectType20509   RE3DataSubjectType = "20509"
	RE3DataSubjectType20510   RE3DataSubjectType = "20510"
	RE3DataSubjectType20511   RE3DataSubjectType = "20511"
	RE3DataSubjectType20512   RE3DataSubjectType = "20512"
	RE3DataSubjectType20513   RE3DataSubjectType = "20513"
	RE3DataSubjectType20514   RE3DataSubjectType = "20514"
	RE3DataSubjectType20515   RE3DataSubjectType = "20515"
	RE3DataSubjectType20516   RE3DataSubjectType = "20516"
	RE3DataSubjectType20517   RE3DataSubjectType = "20517"
	RE3DataSubjectType20518   RE3DataSubjectType = "20518"
	RE3DataSubjectType20519   RE3DataSubjectType = "20519"
	RE3DataSubjectType20520   RE3DataSubjectType = "20520"
	RE3DataSubjectType20521   RE3DataSubjectType = "20521"
	RE3DataSubjectType20522   RE3DataSubjectType = "20522"
	RE3DataSubjectType20523   RE3DataSubjectType = "20523"
	RE3DataSubjectType20524   RE3DataSubjectType = "20524"
	RE3DataSubjectType20525   RE3DataSubjectType = "20525"
	RE3DataSubjectType20526   RE3DataSubjectType = "20526"
	RE3DataSubjectType20527   RE3DataSubjectType = "20527"
	RE3DataSubjectType20528   RE3DataSubjectType = "20528"
	RE3DataSubjectType20529   RE3DataSubjectType = "20529"
	RE3DataSubjectType20530   RE3DataSubjectType = "20530"
	RE3DataSubjectType20531   RE3DataSubjectType = "20531"
	RE3DataSubjectType20532   RE3DataSubjectType = "20532"
	RE3DataSubjectType206     RE3DataSubjectType = "206"
	RE3DataSubjectType20601   RE3DataSubjectType = "20601"
	RE3DataSubjectType20602   RE3DataSubjectType = "20602"
	RE3DataSubjectType20603   RE3DataSubjectType = "20603"
	RE3DataSubjectType20604   RE3DataSubjectType = "20604"
	RE3DataSubjectType20605   RE3DataSubjectType = "20605"
	RE3DataSubjectType20606   RE3DataSubjectType = "20606"
	RE3DataSubjectType20607   RE3DataSubjectType = "20607"
	RE3DataSubjectType20608   RE3DataSubjectType = "20608"
	RE3DataSubjectType20609   RE3DataSubjectType = "20609"
	RE3DataSubjectType20610   RE3DataSubjectType = "20610"
	RE3DataSubjectType20611   RE3DataSubjectType = "20611"
	RE3DataSubjectType23      RE3DataSubjectType = "23"
	RE3DataSubjectType207     RE3DataSubjectType = "207"
	RE3DataSubjectType20701   RE3DataSubjectType = "20701"
	RE3DataSubjectType20702   RE3DataSubjectType = "20702"
	RE3DataSubjectType20703   RE3DataSubjectType = "20703"
	RE3DataSubjectType20704   RE3DataSubjectType = "20704"
	RE3DataSubjectType20705   RE3DataSubjectType = "20705"
	RE3DataSubjectType20706   RE3DataSubjectType = "20706"
	RE3DataSubjectType20707   RE3DataSubjectType = "20707"
	RE3DataSubjectType20708   RE3DataSubjectType = "20708"
	RE3DataSubjectType20709   RE3DataSubjectType = "20709"
	RE3DataSubjectType20710   RE3DataSubjectType = "20710"
	RE3DataSubjectType20711   RE3DataSubjectType = "20711"
	RE3DataSubjectType20712   RE3DataSubjectType = "20712"
	RE3DataSubjectType20713   RE3DataSubjectType = "20713"
	RE3DataSubjectType20714   RE3DataSubjectType = "20714"
	RE3DataSubjectType3       RE3DataSubjectType = "3"
	RE3DataSubjectType31      RE3DataSubjectType = "31"
	RE3DataSubjectType301     RE3DataSubjectType = "301"
	RE3DataSubjectType30101   RE3DataSubjectType = "30101"
	RE3DataSubjectType30102   RE3DataSubjectType = "30102"
	RE3DataSubjectType302     RE3DataSubjectType = "302"
	RE3DataSubjectType30201   RE3DataSubjectType = "30201"
	RE3DataSubjectType30202   RE3DataSubjectType = "30202"
	RE3DataSubjectType30203   RE3DataSubjectType = "30203"
	RE3DataSubjectType303     RE3DataSubjectType = "303"
	RE3DataSubjectType30301   RE3DataSubjectType = "30301"
	RE3DataSubjectType30302   RE3DataSubjectType = "30302"
	RE3DataSubjectType304     RE3DataSubjectType = "304"
	RE3DataSubjectType30401   RE3DataSubjectType = "30401"
	RE3DataSubjectType305     RE3DataSubjectType = "305"
	RE3DataSubjectType30501   RE3DataSubjectType = "30501"
	RE3DataSubjectType30502   RE3DataSubjectType = "30502"
	RE3DataSubjectType306     RE3DataSubjectType = "306"
	RE3DataSubjectType30601   RE3DataSubjectType = "30601"
	RE3DataSubjectType30602   RE3DataSubjectType = "30602"
	RE3DataSubjectType30603   RE3DataSubjectType = "30603"
	RE3DataSubjectType32      RE3DataSubjectType = "32"
	RE3DataSubjectType307     RE3DataSubjectType = "307"
	RE3DataSubjectType30701   RE3DataSubjectType = "30701"
	RE3DataSubjectType30702   RE3DataSubjectType = "30702"
	RE3DataSubjectType308     RE3DataSubjectType = "308"
	RE3DataSubjectType30801   RE3DataSubjectType = "30801"
	RE3DataSubjectType309     RE3DataSubjectType = "309"
	RE3DataSubjectType30901   RE3DataSubjectType = "30901"
	RE3DataSubjectType310     RE3DataSubjectType = "310"
	RE3DataSubjectType31001   RE3DataSubjectType = "31001"
	RE3DataSubjectType311     RE3DataSubjectType = "311"
	RE3DataSubjectType31101   RE3DataSubjectType = "31101"
	RE3DataSubjectType33      RE3DataSubjectType = "33"
	RE3DataSubjectType312     RE3DataSubjectType = "312"
	RE3DataSubjectType31201   RE3DataSubjectType = "31201"
	RE3DataSubjectType34      RE3DataSubjectType = "34"
	RE3DataSubjectType313     RE3DataSubjectType = "313"
	RE3DataSubjectType31301   RE3DataSubjectType = "31301"
	RE3DataSubjectType31302   RE3DataSubjectType = "31302"
	RE3DataSubjectType314     RE3DataSubjectType = "314"
	RE3DataSubjectType31401   RE3DataSubjectType = "31401"
	RE3DataSubjectType315     RE3DataSubjectType = "315"
	RE3DataSubjectType31501   RE3DataSubjectType = "31501"
	RE3DataSubjectType31502   RE3DataSubjectType = "31502"
	RE3DataSubjectType316     RE3DataSubjectType = "316"
	RE3DataSubjectType31601   RE3DataSubjectType = "31601"
	RE3DataSubjectType317     RE3DataSubjectType = "317"
	RE3DataSubjectType31701   RE3DataSubjectType = "31701"
	RE3DataSubjectType31702   RE3DataSubjectType = "31702"
	RE3DataSubjectType318     RE3DataSubjectType = "318"
	RE3DataSubjectType31801   RE3DataSubjectType = "31801"
	RE3DataSubjectType4       RE3DataSubjectType = "4"
	RE3DataSubjectType41      RE3DataSubjectType = "41"
	RE3DataSubjectType401     RE3DataSubjectType = "401"
	RE3DataSubjectType40101   RE3DataSubjectType = "40101"
	RE3DataSubjectType40102   RE3DataSubjectType = "40102"
	RE3DataSubjectType40103   RE3DataSubjectType = "40103"
	RE3DataSubjectType40104   RE3DataSubjectType = "40104"
	RE3DataSubjectType40105   RE3DataSubjectType = "40105"
	RE3DataSubjectType402     RE3DataSubjectType = "402"
	RE3DataSubjectType40201   RE3DataSubjectType = "40201"
	RE3DataSubjectType40202   RE3DataSubjectType = "40202"
	RE3DataSubjectType40203   RE3DataSubjectType = "40203"
	RE3DataSubjectType40204   RE3DataSubjectType = "40204"
	RE3DataSubjectType42      RE3DataSubjectType = "42"
	RE3DataSubjectType403     RE3DataSubjectType = "403"
	RE3DataSubjectType40301   RE3DataSubjectType = "40301"
	RE3DataSubjectType40302   RE3DataSubjectType = "40302"
	RE3DataSubjectType40303   RE3DataSubjectType = "40303"
	RE3DataSubjectType40304   RE3DataSubjectType = "40304"
	RE3DataSubjectType404     RE3DataSubjectType = "404"
	RE3DataSubjectType40401   RE3DataSubjectType = "40401"
	RE3DataSubjectType40402   RE3DataSubjectType = "40402"
	RE3DataSubjectType40403   RE3DataSubjectType = "40403"
	RE3DataSubjectType40404   RE3DataSubjectType = "40404"
	RE3DataSubjectType43      RE3DataSubjectType = "43"
	RE3DataSubjectType405     RE3DataSubjectType = "405"
	RE3DataSubjectType40501   RE3DataSubjectType = "40501"
	RE3DataSubjectType40502   RE3DataSubjectType = "40502"
	RE3DataSubjectType40503   RE3DataSubjectType = "40503"
	RE3DataSubjectType40504   RE3DataSubjectType = "40504"
	RE3DataSubjectType40505   RE3DataSubjectType = "40505"
	RE3DataSubjectType406     RE3DataSubjectType = "406"
	RE3DataSubjectType40601   RE3DataSubjectType = "40601"
	RE3DataSubjectType40602   RE3DataSubjectType = "40602"
	RE3DataSubjectType40603   RE3DataSubjectType = "40603"
	RE3DataSubjectType40604   RE3DataSubjectType = "40604"
	RE3DataSubjectType40605   RE3DataSubjectType = "40605"
	RE3DataSubjectType44      RE3DataSubjectType = "44"
	RE3DataSubjectType407     RE3DataSubjectType = "407"
	RE3DataSubjectType40701   RE3DataSubjectType = "40701"
	RE3DataSubjectType40702   RE3DataSubjectType = "40702"
	RE3DataSubjectType40703   RE3DataSubjectType = "40703"
	RE3DataSubjectType40704   RE3DataSubjectType = "40704"
	RE3DataSubjectType40705   RE3DataSubjectType = "40705"
	RE3DataSubjectType408     RE3DataSubjectType = "408"
	RE3DataSubjectType40801   RE3DataSubjectType = "40801"
	RE3DataSubjectType40802   RE3DataSubjectType = "40802"
	RE3DataSubjectType40803   RE3DataSubjectType = "40803"
	RE3DataSubjectType409     RE3DataSubjectType = "409"
	RE3DataSubjectType40901   RE3DataSubjectType = "40901"
	RE3DataSubjectType40902   RE3DataSubjectType = "40902"
	RE3DataSubjectType40903   RE3DataSubjectType = "40903"
	RE3DataSubjectType40904   RE3DataSubjectType = "40904"
	RE3DataSubjectType40905   RE3DataSubjectType = "40905"
	RE3DataSubjectType45      RE3DataSubjectType = "45"
	RE3DataSubjectType410     RE3DataSubjectType = "410"
	RE3DataSubjectType41001   RE3DataSubjectType = "41001"
	RE3DataSubjectType41002   RE3DataSubjectType = "41002"
	RE3DataSubjectType41003   RE3DataSubjectType = "41003"
	RE3DataSubjectType41004   RE3DataSubjectType = "41004"
	RE3DataSubjectType41005   RE3DataSubjectType = "41005"
	RE3DataSubjectType41006   RE3DataSubjectType = "41006"
)

var RE3DataSubjectTypeReverse = map[string]RE3DataSubjectType{
	string(RE3DataSubjectType1):       RE3DataSubjectType1,
	string(RE3DataSubjectType11):      RE3DataSubjectType11,
	string(RE3DataSubjectType101):     RE3DataSubjectType101,
	string(RE3DataSubjectType10101):   RE3DataSubjectType10101,
	string(RE3DataSubjectType10102):   RE3DataSubjectType10102,
	string(RE3DataSubjectType10103):   RE3DataSubjectType10103,
	string(RE3DataSubjectType10104):   RE3DataSubjectType10104,
	string(RE3DataSubjectType10105):   RE3DataSubjectType10105,
	string(RE3DataSubjectType102):     RE3DataSubjectType102,
	string(RE3DataSubjectType10201):   RE3DataSubjectType10201,
	string(RE3DataSubjectType10202):   RE3DataSubjectType10202,
	string(RE3DataSubjectType10203):   RE3DataSubjectType10203,
	string(RE3DataSubjectType10204):   RE3DataSubjectType10204,
	string(RE3DataSubjectType103):     RE3DataSubjectType103,
	string(RE3DataSubjectType10301):   RE3DataSubjectType10301,
	string(RE3DataSubjectType10302):   RE3DataSubjectType10302,
	string(RE3DataSubjectType10303):   RE3DataSubjectType10303,
	string(RE3DataSubjectType104):     RE3DataSubjectType104,
	string(RE3DataSubjectType10401):   RE3DataSubjectType10401,
	string(RE3DataSubjectType10402):   RE3DataSubjectType10402,
	string(RE3DataSubjectType10403):   RE3DataSubjectType10403,
	string(RE3DataSubjectType105):     RE3DataSubjectType105,
	string(RE3DataSubjectType10501):   RE3DataSubjectType10501,
	string(RE3DataSubjectType10502):   RE3DataSubjectType10502,
	string(RE3DataSubjectType10503):   RE3DataSubjectType10503,
	string(RE3DataSubjectType10504):   RE3DataSubjectType10504,
	string(RE3DataSubjectType106):     RE3DataSubjectType106,
	string(RE3DataSubjectType10601):   RE3DataSubjectType10601,
	string(RE3DataSubjectType10602):   RE3DataSubjectType10602,
	string(RE3DataSubjectType10603):   RE3DataSubjectType10603,
	string(RE3DataSubjectType10604):   RE3DataSubjectType10604,
	string(RE3DataSubjectType10605):   RE3DataSubjectType10605,
	string(RE3DataSubjectType107):     RE3DataSubjectType107,
	string(RE3DataSubjectType10701):   RE3DataSubjectType10701,
	string(RE3DataSubjectType10702):   RE3DataSubjectType10702,
	string(RE3DataSubjectType108):     RE3DataSubjectType108,
	string(RE3DataSubjectType10801):   RE3DataSubjectType10801,
	string(RE3DataSubjectType10802):   RE3DataSubjectType10802,
	string(RE3DataSubjectType10803):   RE3DataSubjectType10803,
	string(RE3DataSubjectType12):      RE3DataSubjectType12,
	string(RE3DataSubjectType109):     RE3DataSubjectType109,
	string(RE3DataSubjectType10901):   RE3DataSubjectType10901,
	string(RE3DataSubjectType10902):   RE3DataSubjectType10902,
	string(RE3DataSubjectType10903):   RE3DataSubjectType10903,
	string(RE3DataSubjectType110):     RE3DataSubjectType110,
	string(RE3DataSubjectType11001):   RE3DataSubjectType11001,
	string(RE3DataSubjectType11002):   RE3DataSubjectType11002,
	string(RE3DataSubjectType11003):   RE3DataSubjectType11003,
	string(RE3DataSubjectType11004):   RE3DataSubjectType11004,
	string(RE3DataSubjectType111):     RE3DataSubjectType111,
	string(RE3DataSubjectType11101):   RE3DataSubjectType11101,
	string(RE3DataSubjectType11102):   RE3DataSubjectType11102,
	string(RE3DataSubjectType11103):   RE3DataSubjectType11103,
	string(RE3DataSubjectType11104):   RE3DataSubjectType11104,
	string(RE3DataSubjectType112):     RE3DataSubjectType112,
	string(RE3DataSubjectType11201):   RE3DataSubjectType11201,
	string(RE3DataSubjectType11202):   RE3DataSubjectType11202,
	string(RE3DataSubjectType11203):   RE3DataSubjectType11203,
	string(RE3DataSubjectType11204):   RE3DataSubjectType11204,
	string(RE3DataSubjectType11205):   RE3DataSubjectType11205,
	string(RE3DataSubjectType11206y):  RE3DataSubjectType11206y,
	string(RE3DataSubjectType113):     RE3DataSubjectType113,
	string(RE3DataSubjectType11301):   RE3DataSubjectType11301,
	string(RE3DataSubjectType11302):   RE3DataSubjectType11302,
	string(RE3DataSubjectType11303):   RE3DataSubjectType11303,
	string(RE3DataSubjectType11304):   RE3DataSubjectType11304,
	string(RE3DataSubjectType11305):   RE3DataSubjectType11305,
	string(RE3DataSubjectType2):       RE3DataSubjectType2,
	string(RE3DataSubjectType21):      RE3DataSubjectType21,
	string(RE3DataSubjectType201):     RE3DataSubjectType201,
	string(RE3DataSubjectType20101):   RE3DataSubjectType20101,
	string(RE3DataSubjectType20102):   RE3DataSubjectType20102,
	string(RE3DataSubjectType20103):   RE3DataSubjectType20103,
	string(RE3DataSubjectType20104):   RE3DataSubjectType20104,
	string(RE3DataSubjectType20105):   RE3DataSubjectType20105,
	string(RE3DataSubjectType20106):   RE3DataSubjectType20106,
	string(RE3DataSubjectType20107):   RE3DataSubjectType20107,
	string(RE3DataSubjectType20108):   RE3DataSubjectType20108,
	string(RE3DataSubjectType202):     RE3DataSubjectType202,
	string(RE3DataSubjectType20201):   RE3DataSubjectType20201,
	string(RE3DataSubjectType20202):   RE3DataSubjectType20202,
	string(RE3DataSubjectType20203):   RE3DataSubjectType20203,
	string(RE3DataSubjectType20204):   RE3DataSubjectType20204,
	string(RE3DataSubjectType20205):   RE3DataSubjectType20205,
	string(RE3DataSubjectType20206):   RE3DataSubjectType20206,
	string(RE3DataSubjectType20207):   RE3DataSubjectType20207,
	string(RE3DataSubjectType203):     RE3DataSubjectType203,
	string(RE3DataSubjectType20301):   RE3DataSubjectType20301,
	string(RE3DataSubjectType20302):   RE3DataSubjectType20302,
	string(RE3DataSubjectType20303):   RE3DataSubjectType20303,
	string(RE3DataSubjectType20304):   RE3DataSubjectType20304,
	string(RE3DataSubjectType20305):   RE3DataSubjectType20305,
	string(RE3DataSubjectType20306):   RE3DataSubjectType20306,
	string(RE3DataSubjectType22):      RE3DataSubjectType22,
	string(RE3DataSubjectType204):     RE3DataSubjectType204,
	string(RE3DataSubjectType20401):   RE3DataSubjectType20401,
	string(RE3DataSubjectType20402):   RE3DataSubjectType20402,
	string(RE3DataSubjectType20403):   RE3DataSubjectType20403,
	string(RE3DataSubjectType20404):   RE3DataSubjectType20404,
	string(RE3DataSubjectType20405):   RE3DataSubjectType20405,
	string(RE3DataSubjectType205):     RE3DataSubjectType205,
	string(RE3DataSubjectType20501cs): RE3DataSubjectType20501cs,
	string(RE3DataSubjectType20502):   RE3DataSubjectType20502,
	string(RE3DataSubjectType20503):   RE3DataSubjectType20503,
	string(RE3DataSubjectType20504):   RE3DataSubjectType20504,
	string(RE3DataSubjectType20505):   RE3DataSubjectType20505,
	string(RE3DataSubjectType20506):   RE3DataSubjectType20506,
	string(RE3DataSubjectType20507):   RE3DataSubjectType20507,
	string(RE3DataSubjectType20508):   RE3DataSubjectType20508,
	string(RE3DataSubjectType20509):   RE3DataSubjectType20509,
	string(RE3DataSubjectType20510):   RE3DataSubjectType20510,
	string(RE3DataSubjectType20511):   RE3DataSubjectType20511,
	string(RE3DataSubjectType20512):   RE3DataSubjectType20512,
	string(RE3DataSubjectType20513):   RE3DataSubjectType20513,
	string(RE3DataSubjectType20514):   RE3DataSubjectType20514,
	string(RE3DataSubjectType20515):   RE3DataSubjectType20515,
	string(RE3DataSubjectType20516):   RE3DataSubjectType20516,
	string(RE3DataSubjectType20517):   RE3DataSubjectType20517,
	string(RE3DataSubjectType20518):   RE3DataSubjectType20518,
	string(RE3DataSubjectType20519):   RE3DataSubjectType20519,
	string(RE3DataSubjectType20520):   RE3DataSubjectType20520,
	string(RE3DataSubjectType20521):   RE3DataSubjectType20521,
	string(RE3DataSubjectType20522):   RE3DataSubjectType20522,
	string(RE3DataSubjectType20523):   RE3DataSubjectType20523,
	string(RE3DataSubjectType20524):   RE3DataSubjectType20524,
	string(RE3DataSubjectType20525):   RE3DataSubjectType20525,
	string(RE3DataSubjectType20526):   RE3DataSubjectType20526,
	string(RE3DataSubjectType20527):   RE3DataSubjectType20527,
	string(RE3DataSubjectType20528):   RE3DataSubjectType20528,
	string(RE3DataSubjectType20529):   RE3DataSubjectType20529,
	string(RE3DataSubjectType20530):   RE3DataSubjectType20530,
	string(RE3DataSubjectType20531):   RE3DataSubjectType20531,
	string(RE3DataSubjectType20532):   RE3DataSubjectType20532,
	string(RE3DataSubjectType206):     RE3DataSubjectType206,
	string(RE3DataSubjectType20601):   RE3DataSubjectType20601,
	string(RE3DataSubjectType20602):   RE3DataSubjectType20602,
	string(RE3DataSubjectType20603):   RE3DataSubjectType20603,
	string(RE3DataSubjectType20604):   RE3DataSubjectType20604,
	string(RE3DataSubjectType20605):   RE3DataSubjectType20605,
	string(RE3DataSubjectType20606):   RE3DataSubjectType20606,
	string(RE3DataSubjectType20607):   RE3DataSubjectType20607,
	string(RE3DataSubjectType20608):   RE3DataSubjectType20608,
	string(RE3DataSubjectType20609):   RE3DataSubjectType20609,
	string(RE3DataSubjectType20610):   RE3DataSubjectType20610,
	string(RE3DataSubjectType20611):   RE3DataSubjectType20611,
	string(RE3DataSubjectType23):      RE3DataSubjectType23,
	string(RE3DataSubjectType207):     RE3DataSubjectType207,
	string(RE3DataSubjectType20701):   RE3DataSubjectType20701,
	string(RE3DataSubjectType20702):   RE3DataSubjectType20702,
	string(RE3DataSubjectType20703):   RE3DataSubjectType20703,
	string(RE3DataSubjectType20704):   RE3DataSubjectType20704,
	string(RE3DataSubjectType20705):   RE3DataSubjectType20705,
	string(RE3DataSubjectType20706):   RE3DataSubjectType20706,
	string(RE3DataSubjectType20707):   RE3DataSubjectType20707,
	string(RE3DataSubjectType20708):   RE3DataSubjectType20708,
	string(RE3DataSubjectType20709):   RE3DataSubjectType20709,
	string(RE3DataSubjectType20710):   RE3DataSubjectType20710,
	string(RE3DataSubjectType20711):   RE3DataSubjectType20711,
	string(RE3DataSubjectType20712):   RE3DataSubjectType20712,
	string(RE3DataSubjectType20713):   RE3DataSubjectType20713,
	string(RE3DataSubjectType20714):   RE3DataSubjectType20714,
	string(RE3DataSubjectType3):       RE3DataSubjectType3,
	string(RE3DataSubjectType31):      RE3DataSubjectType31,
	string(RE3DataSubjectType301):     RE3DataSubjectType301,
	string(RE3DataSubjectType30101):   RE3DataSubjectType30101,
	string(RE3DataSubjectType30102):   RE3DataSubjectType30102,
	string(RE3DataSubjectType302):     RE3DataSubjectType302,
	string(RE3DataSubjectType30201):   RE3DataSubjectType30201,
	string(RE3DataSubjectType30202):   RE3DataSubjectType30202,
	string(RE3DataSubjectType30203):   RE3DataSubjectType30203,
	string(RE3DataSubjectType303):     RE3DataSubjectType303,
	string(RE3DataSubjectType30301):   RE3DataSubjectType30301,
	string(RE3DataSubjectType30302):   RE3DataSubjectType30302,
	string(RE3DataSubjectType304):     RE3DataSubjectType304,
	string(RE3DataSubjectType30401):   RE3DataSubjectType30401,
	string(RE3DataSubjectType305):     RE3DataSubjectType305,
	string(RE3DataSubjectType30501):   RE3DataSubjectType30501,
	string(RE3DataSubjectType30502):   RE3DataSubjectType30502,
	string(RE3DataSubjectType306):     RE3DataSubjectType306,
	string(RE3DataSubjectType30601):   RE3DataSubjectType30601,
	string(RE3DataSubjectType30602):   RE3DataSubjectType30602,
	string(RE3DataSubjectType30603):   RE3DataSubjectType30603,
	string(RE3DataSubjectType32):      RE3DataSubjectType32,
	string(RE3DataSubjectType307):     RE3DataSubjectType307,
	string(RE3DataSubjectType30701):   RE3DataSubjectType30701,
	string(RE3DataSubjectType30702):   RE3DataSubjectType30702,
	string(RE3DataSubjectType308):     RE3DataSubjectType308,
	string(RE3DataSubjectType30801):   RE3DataSubjectType30801,
	string(RE3DataSubjectType309):     RE3DataSubjectType309,
	string(RE3DataSubjectType30901):   RE3DataSubjectType30901,
	string(RE3DataSubjectType310):     RE3DataSubjectType310,
	string(RE3DataSubjectType31001):   RE3DataSubjectType31001,
	string(RE3DataSubjectType311):     RE3DataSubjectType311,
	string(RE3DataSubjectType31101):   RE3DataSubjectType31101,
	string(RE3DataSubjectType33):      RE3DataSubjectType33,
	string(RE3DataSubjectType312):     RE3DataSubjectType312,
	string(RE3DataSubjectType31201):   RE3DataSubjectType31201,
	string(RE3DataSubjectType34):      RE3DataSubjectType34,
	string(RE3DataSubjectType313):     RE3DataSubjectType313,
	string(RE3DataSubjectType31301):   RE3DataSubjectType31301,
	string(RE3DataSubjectType31302):   RE3DataSubjectType31302,
	string(RE3DataSubjectType314):     RE3DataSubjectType314,
	string(RE3DataSubjectType31401):   RE3DataSubjectType31401,
	string(RE3DataSubjectType315):     RE3DataSubjectType315,
	string(RE3DataSubjectType31501):   RE3DataSubjectType31501,
	string(RE3DataSubjectType31502):   RE3DataSubjectType31502,
	string(RE3DataSubjectType316):     RE3DataSubjectType316,
	string(RE3DataSubjectType31601):   RE3DataSubjectType31601,
	string(RE3DataSubjectType317):     RE3DataSubjectType317,
	string(RE3DataSubjectType31701):   RE3DataSubjectType31701,
	string(RE3DataSubjectType31702):   RE3DataSubjectType31702,
	string(RE3DataSubjectType318):     RE3DataSubjectType318,
	string(RE3DataSubjectType31801):   RE3DataSubjectType31801,
	string(RE3DataSubjectType4):       RE3DataSubjectType4,
	string(RE3DataSubjectType41):      RE3DataSubjectType41,
	string(RE3DataSubjectType401):     RE3DataSubjectType401,
	string(RE3DataSubjectType40101):   RE3DataSubjectType40101,
	string(RE3DataSubjectType40102):   RE3DataSubjectType40102,
	string(RE3DataSubjectType40103):   RE3DataSubjectType40103,
	string(RE3DataSubjectType40104):   RE3DataSubjectType40104,
	string(RE3DataSubjectType40105):   RE3DataSubjectType40105,
	string(RE3DataSubjectType402):     RE3DataSubjectType402,
	string(RE3DataSubjectType40201):   RE3DataSubjectType40201,
	string(RE3DataSubjectType40202):   RE3DataSubjectType40202,
	string(RE3DataSubjectType40203):   RE3DataSubjectType40203,
	string(RE3DataSubjectType40204):   RE3DataSubjectType40204,
	string(RE3DataSubjectType42):      RE3DataSubjectType42,
	string(RE3DataSubjectType403):     RE3DataSubjectType403,
	string(RE3DataSubjectType40301):   RE3DataSubjectType40301,
	string(RE3DataSubjectType40302):   RE3DataSubjectType40302,
	string(RE3DataSubjectType40303):   RE3DataSubjectType40303,
	string(RE3DataSubjectType40304):   RE3DataSubjectType40304,
	string(RE3DataSubjectType404):     RE3DataSubjectType404,
	string(RE3DataSubjectType40401):   RE3DataSubjectType40401,
	string(RE3DataSubjectType40402):   RE3DataSubjectType40402,
	string(RE3DataSubjectType40403):   RE3DataSubjectType40403,
	string(RE3DataSubjectType40404):   RE3DataSubjectType40404,
	string(RE3DataSubjectType43):      RE3DataSubjectType43,
	string(RE3DataSubjectType405):     RE3DataSubjectType405,
	string(RE3DataSubjectType40501):   RE3DataSubjectType40501,
	string(RE3DataSubjectType40502):   RE3DataSubjectType40502,
	string(RE3DataSubjectType40503):   RE3DataSubjectType40503,
	string(RE3DataSubjectType40504):   RE3DataSubjectType40504,
	string(RE3DataSubjectType40505):   RE3DataSubjectType40505,
	string(RE3DataSubjectType406):     RE3DataSubjectType406,
	string(RE3DataSubjectType40601):   RE3DataSubjectType40601,
	string(RE3DataSubjectType40602):   RE3DataSubjectType40602,
	string(RE3DataSubjectType40603):   RE3DataSubjectType40603,
	string(RE3DataSubjectType40604):   RE3DataSubjectType40604,
	string(RE3DataSubjectType40605):   RE3DataSubjectType40605,
	string(RE3DataSubjectType44):      RE3DataSubjectType44,
	string(RE3DataSubjectType407):     RE3DataSubjectType407,
	string(RE3DataSubjectType40701):   RE3DataSubjectType40701,
	string(RE3DataSubjectType40702):   RE3DataSubjectType40702,
	string(RE3DataSubjectType40703):   RE3DataSubjectType40703,
	string(RE3DataSubjectType40704):   RE3DataSubjectType40704,
	string(RE3DataSubjectType40705):   RE3DataSubjectType40705,
	string(RE3DataSubjectType408):     RE3DataSubjectType408,
	string(RE3DataSubjectType40801):   RE3DataSubjectType40801,
	string(RE3DataSubjectType40802):   RE3DataSubjectType40802,
	string(RE3DataSubjectType40803):   RE3DataSubjectType40803,
	string(RE3DataSubjectType409):     RE3DataSubjectType409,
	string(RE3DataSubjectType40901):   RE3DataSubjectType40901,
	string(RE3DataSubjectType40902):   RE3DataSubjectType40902,
	string(RE3DataSubjectType40903):   RE3DataSubjectType40903,
	string(RE3DataSubjectType40904):   RE3DataSubjectType40904,
	string(RE3DataSubjectType40905):   RE3DataSubjectType40905,
	string(RE3DataSubjectType45):      RE3DataSubjectType45,
	string(RE3DataSubjectType410):     RE3DataSubjectType410,
	string(RE3DataSubjectType41001):   RE3DataSubjectType41001,
	string(RE3DataSubjectType41002):   RE3DataSubjectType41002,
	string(RE3DataSubjectType41003):   RE3DataSubjectType41003,
	string(RE3DataSubjectType41004):   RE3DataSubjectType41004,
	string(RE3DataSubjectType41005):   RE3DataSubjectType41005,
	string(RE3DataSubjectType41006):   RE3DataSubjectType41006,
}

var RE3DataSubjectTypeName = map[RE3DataSubjectType]string{
	RE3DataSubjectType1:       "Humanities and Social Sciences",
	RE3DataSubjectType11:      "Humanities",
	RE3DataSubjectType101:     "Ancient Cultures",
	RE3DataSubjectType10101:   "Prehistory",
	RE3DataSubjectType10102:   "Classical Philology",
	RE3DataSubjectType10103:   "Ancient History",
	RE3DataSubjectType10104:   "Classical Archaeology",
	RE3DataSubjectType10105:   "Egyptology and Ancient Near Eastern Studies",
	RE3DataSubjectType102:     "History",
	RE3DataSubjectType10201:   "Medieval History",
	RE3DataSubjectType10202:   "Early Modern History",
	RE3DataSubjectType10203:   "Modern and Current History",
	RE3DataSubjectType10204:   "History of Science",
	RE3DataSubjectType103:     "Fine Arts, Music, Theatre and Media Studies",
	RE3DataSubjectType10301:   "Art History",
	RE3DataSubjectType10302:   "Musicology",
	RE3DataSubjectType10303:   "Theatre and Media Studies",
	RE3DataSubjectType104:     "Linguistics",
	RE3DataSubjectType10401:   "General and Applied Linguistics",
	RE3DataSubjectType10402:   "Individual Linguistics",
	RE3DataSubjectType10403:   "Typology, Non-European Languages, Historical Linguistics",
	RE3DataSubjectType105:     "Literary Studies",
	RE3DataSubjectType10501:   "Medieval German Literature",
	RE3DataSubjectType10502:   "Modern German Literature",
	RE3DataSubjectType10503:   "European and American Literature",
	RE3DataSubjectType10504:   "General and Comparative Literature and Cultural Studies",
	RE3DataSubjectType106:     "Non-European Languages and Cultures, Social and Cultural Anthropology, Jewish Studies and Religious Studies",
	RE3DataSubjectType10601:   "Social and Cultural Anthropology and Ethnology/Folklore",
	RE3DataSubjectType10602:   "Asian Studies",
	RE3DataSubjectType10603:   "African, American and Oceania Studies",
	RE3DataSubjectType10604:   "Islamic Studies, Arabian Studies, Semitic Studies",
	RE3DataSubjectType10605:   "Religious Studies and Jewish Studies",
	RE3DataSubjectType107:     "Theology",
	RE3DataSubjectType10701:   "Protestant Theology",
	RE3DataSubjectType10702:   "Roman Catholic Theology",
	RE3DataSubjectType108:     "Philosophy",
	RE3DataSubjectType10801:   "History of Philosophy",
	RE3DataSubjectType10802:   "Theoretical Philosophy",
	RE3DataSubjectType10803:   "Practical Philosophy",
	RE3DataSubjectType12:      "Social and Behavioural Sciences",
	RE3DataSubjectType109:     "Education Sciences",
	RE3DataSubjectType10901:   "General Education and History of Education",
	RE3DataSubjectType10902:   "Research on Teaching, Learning and Training",
	RE3DataSubjectType10903:   "Research on Socialization and Educational Institutions and Professions",
	RE3DataSubjectType110:     "Psychology",
	RE3DataSubjectType11001:   "General, Biological and Mathematical Psychology",
	RE3DataSubjectType11002:   "Developmental and Educational Psychology",
	RE3DataSubjectType11003:   "Social Psychology, Industrial and Organisational Psychology",
	RE3DataSubjectType11004:   "Differential Psychology, Clinical Psychology, Medical Psychology, Methodology",
	RE3DataSubjectType111:     "Social Sciences",
	RE3DataSubjectType11101:   "Sociological Theory",
	RE3DataSubjectType11102:   "Empirical Social Research",
	RE3DataSubjectType11103:   "Communication Science",
	RE3DataSubjectType11104:   "Political Science",
	RE3DataSubjectType112:     "Economics",
	RE3DataSubjectType11201:   "Economic Theory",
	RE3DataSubjectType11202:   "Economic and Social Policy",
	RE3DataSubjectType11203:   "Public Finance",
	RE3DataSubjectType11204:   "Business Administration",
	RE3DataSubjectType11205:   "Statistics and Econometrics",
	RE3DataSubjectType11206y:  "Economic and Social History",
	RE3DataSubjectType113:     "Jurisprudence",
	RE3DataSubjectType11301:   "Legal and Political Philosophy, Legal History, Legal Theory",
	RE3DataSubjectType11302:   "Private Law",
	RE3DataSubjectType11303:   "Public Law",
	RE3DataSubjectType11304:   "Criminal Law and Law of Criminal Procedure",
	RE3DataSubjectType11305:   "Criminology",
	RE3DataSubjectType2:       "Life Sciences",
	RE3DataSubjectType21:      "Biology",
	RE3DataSubjectType201:     "Basic Biological and Medical Research",
	RE3DataSubjectType20101:   "Biochemistry",
	RE3DataSubjectType20102:   "Biophysics",
	RE3DataSubjectType20103:   "Cell Biology",
	RE3DataSubjectType20104:   "Structural Biology",
	RE3DataSubjectType20105:   "General Genetics",
	RE3DataSubjectType20106:   "Developmental Biology",
	RE3DataSubjectType20107:   "Bioinformatics and Theoretical Biology",
	RE3DataSubjectType20108:   "Anatomy",
	RE3DataSubjectType202:     "Plant Sciences",
	RE3DataSubjectType20201:   "Plant Systematics and Evolution",
	RE3DataSubjectType20202:   "Plant Ecology and Ecosystem Analysis",
	RE3DataSubjectType20203:   "Inter-organismic Interactions of Plants",
	RE3DataSubjectType20204:   "Plant Physiology",
	RE3DataSubjectType20205:   "Plant Biochemistry and Biophysics",
	RE3DataSubjectType20206:   "Plant Cell and Developmental Biology",
	RE3DataSubjectType20207:   "Plant Genetics",
	RE3DataSubjectType203:     "Zoology",
	RE3DataSubjectType20301:   "Systematics and Morphology",
	RE3DataSubjectType20302:   "Evolution, Anthropology",
	RE3DataSubjectType20303:   "Animal Ecology, Biodiversity and Ecosystem Research",
	RE3DataSubjectType20304:   "Sensory and Behavioural Biology",
	RE3DataSubjectType20305:   "Biochemistry and Animal Physiology",
	RE3DataSubjectType20306:   "Animal Genetics, Cell and Developmental Biology",
	RE3DataSubjectType22:      "Medicine",
	RE3DataSubjectType204:     "Microbiology, Virology and Immunology",
	RE3DataSubjectType20401:   "Metabolism, Biochemistry and Genetics of Microorganisms",
	RE3DataSubjectType20402:   "Microbial Ecology and Applied Microbiology",
	RE3DataSubjectType20403:   "Medical Microbiology, Molecular Infection Biology",
	RE3DataSubjectType20404:   "Virology",
	RE3DataSubjectType20405:   "Immunology",
	RE3DataSubjectType205:     "Medicine",
	RE3DataSubjectType20501cs: "Epidemiology, Medical Biometry, Medical Informatics",
	RE3DataSubjectType20502:   "Public Health, Health Services Research, Social Medicine",
	RE3DataSubjectType20503:   "Human Genetics",
	RE3DataSubjectType20504:   "Physiology",
	RE3DataSubjectType20505:   "Nutritional Sciences",
	RE3DataSubjectType20506:   "Pathology and Forensic Medicine",
	RE3DataSubjectType20507:   "Clinical Chemistry and Pathobiochemistry",
	RE3DataSubjectType20508:   "Pharmacy",
	RE3DataSubjectType20509:   "Pharmacology",
	RE3DataSubjectType20510:   "Toxicology and Occupational Medicine",
	RE3DataSubjectType20511:   "Anaesthesiology",
	RE3DataSubjectType20512:   "Cardiology, Angiology",
	RE3DataSubjectType20513:   "Pneumology, Clinical Infectiology Intensive Care Medicine",
	RE3DataSubjectType20514:   "Hematology, Oncology, Transfusion Medicine",
	RE3DataSubjectType20515:   "Gastroenterology, Metabolism",
	RE3DataSubjectType20516:   "Nephrology",
	RE3DataSubjectType20517:   "Endocrinology, Diabetology",
	RE3DataSubjectType20518:   "Rheumatology, Clinical Immunology, Allergology",
	RE3DataSubjectType20519:   "Dermatology",
	RE3DataSubjectType20520:   "Pediatric and Adolescent Medicine",
	RE3DataSubjectType20521:   "Gynaecology and Obstetrics",
	RE3DataSubjectType20522:   "Reproductive Medicine/Biology",
	RE3DataSubjectType20523:   "Urology",
	RE3DataSubjectType20524:   "Gerontology and Geriatric Medicine",
	RE3DataSubjectType20525:   "Vascular and Visceral Surgery",
	RE3DataSubjectType20526:   "Cardiothoracic Surgery",
	RE3DataSubjectType20527:   "Traumatology and Orthopaedics",
	RE3DataSubjectType20528:   "Dentistry, Oral Surgery",
	RE3DataSubjectType20529:   "Otolaryngology",
	RE3DataSubjectType20530:   "Radiology and Nuclear Medicine",
	RE3DataSubjectType20531:   "Radiation Oncology and Radiobiology",
	RE3DataSubjectType20532:   "Biomedical Technology and Medical Physics",
	RE3DataSubjectType206:     "Neurosciences",
	RE3DataSubjectType20601:   "Molecular Neuroscience and Neurogenetics",
	RE3DataSubjectType20602:   "Cellular Neuroscience",
	RE3DataSubjectType20603:   "Developmental Neurobiology",
	RE3DataSubjectType20604:   "Systemic Neuroscience, Computational Neuroscience, Behaviour",
	RE3DataSubjectType20605:   "Comparative Neurobiology",
	RE3DataSubjectType20606:   "Cognitive Neuroscience and Neuroimaging",
	RE3DataSubjectType20607:   "Molecular Neurology",
	RE3DataSubjectType20608:   "Clinical Neurosciences I - Neurology, Neurosurgery",
	RE3DataSubjectType20609:   "Biological Psychiatry",
	RE3DataSubjectType20610:   "Clinical Neurosciences II - Psychiatry, Psychotherapy, Psychosomatic Medicine",
	RE3DataSubjectType20611:   "Clinical Neurosciences III - Ophthalmology",
	RE3DataSubjectType23:      "Agriculture, Forestry, Horticulture and Veterinary Medicine",
	RE3DataSubjectType207:     "Agriculture, Forestry, Horticulture and Veterinary Medicine",
	RE3DataSubjectType20701:   "Soil Sciences",
	RE3DataSubjectType20702:   "Plant Cultivation",
	RE3DataSubjectType20703:   "Plant Nutrition",
	RE3DataSubjectType20704:   "Ecology of Agricultural Landscapes",
	RE3DataSubjectType20705:   "Plant Breeding",
	RE3DataSubjectType20706:   "Phytomedicine",
	RE3DataSubjectType20707:   "Agricultural and Food Process Engineering",
	RE3DataSubjectType20708:   "Agricultural Economics and Sociology",
	RE3DataSubjectType20709:   "Inventory Control and Use of Forest Resources",
	RE3DataSubjectType20710:   "Basic Forest Research",
	RE3DataSubjectType20711:   "Animal Husbandry, Breeding and Hygiene",
	RE3DataSubjectType20712:   "Animal Nutrition and Nutrition Physiology",
	RE3DataSubjectType20713:   "Basic Veterinary Medical Science",
	RE3DataSubjectType20714:   "Basic Research on Pathogenesis, Diagnostics and Therapy and Clinical Veterinary Medicine",
	RE3DataSubjectType3:       "Natural Sciences",
	RE3DataSubjectType31:      "Chemistry",
	RE3DataSubjectType301:     "Molecular Chemistry",
	RE3DataSubjectType30101:   "Inorganic Molecular Chemistry",
	RE3DataSubjectType30102:   "Organic Molecular Chemistry",
	RE3DataSubjectType302:     "Chemical Solid State and Surface Research",
	RE3DataSubjectType30201:   "Solid State and Surface Chemistry, Material Synthesis",
	RE3DataSubjectType30202:   "Physical Chemistry of Solids and Surfaces, Material Characterisation",
	RE3DataSubjectType30203:   "Theory and Modelling",
	RE3DataSubjectType303:     "Physical and Theoretical Chemistry",
	RE3DataSubjectType30301:   "Physical Chemistry of Molecules, Interfaces and Liquids - Spectroscopy, Kinetics",
	RE3DataSubjectType30302:   "General Theoretical Chemistry",
	RE3DataSubjectType304:     "Analytical Chemistry, Method Development (Chemistry)",
	RE3DataSubjectType30401:   "Analytical Chemistry, Method Development (Chemistry)",
	RE3DataSubjectType305:     "Biological Chemistry and Food Chemistry",
	RE3DataSubjectType30501:   "Biological and Biomimetic Chemistry",
	RE3DataSubjectType30502:   "Food Chemistry",
	RE3DataSubjectType306:     "Polymer Research",
	RE3DataSubjectType30601:   "Preparatory and Physical Chemistry of Polymers",
	RE3DataSubjectType30602:   "Experimental and Theoretical Physics of Polymers",
	RE3DataSubjectType30603:   "Polymer Materials",
	RE3DataSubjectType32:      "Physics",
	RE3DataSubjectType307:     "Condensed Matter Physics",
	RE3DataSubjectType30701:   "Experimental Condensed Matter Physics",
	RE3DataSubjectType30702:   "Theoretical Condensed Matter Physics",
	RE3DataSubjectType308:     "Optics, Quantum Optics and Physics of Atoms, Molecules and Plasmas",
	RE3DataSubjectType30801:   "Optics, Quantum Optics, Atoms, Molecules, Plasmas",
	RE3DataSubjectType309:     "Particles, Nuclei and Fields",
	RE3DataSubjectType30901:   "Particles, Nuclei and Fields",
	RE3DataSubjectType310:     "Statistical Physics, Soft Matter, Biological Physics, Nonlinear Dynamics",
	RE3DataSubjectType31001:   "Statistical Physics, Soft Matter, Biological Physics, Nonlinear Dynamics",
	RE3DataSubjectType311:     "Astrophysics and Astronomy",
	RE3DataSubjectType31101:   "Astrophysics and Astronomy",
	RE3DataSubjectType33:      "Mathematics",
	RE3DataSubjectType312:     "Mathematics",
	RE3DataSubjectType31201:   "Mathematics",
	RE3DataSubjectType34:      "Geosciences (including Geography)",
	RE3DataSubjectType313:     "Atmospheric Science and Oceanography",
	RE3DataSubjectType31301:   "Atmospheric Science",
	RE3DataSubjectType31302:   "Oceanography",
	RE3DataSubjectType314:     "Geology and Palaeontology",
	RE3DataSubjectType31401:   "Geology and Palaeontology",
	RE3DataSubjectType315:     "Geophysics and Geodesy",
	RE3DataSubjectType31501:   "Geophysics",
	RE3DataSubjectType31502:   "Geodesy, Photogrammetry, Remote Sensing, Geoinformatics, Cartogaphy",
	RE3DataSubjectType316:     "Geochemistry, Mineralogy and Crystallography",
	RE3DataSubjectType31601:   "Geochemistry, Mineralogy and Crystallography",
	RE3DataSubjectType317:     "Geography",
	RE3DataSubjectType31701:   "Physical Geography",
	RE3DataSubjectType31702:   "Human Geography",
	RE3DataSubjectType318:     "Water Research",
	RE3DataSubjectType31801:   "Hydrogeology, Hydrology, Limnology, Urban Water Management, Water Chemistry, Integrated Water Resources Management",
	RE3DataSubjectType4:       "Engineering Sciences",
	RE3DataSubjectType41:      "Mechanical and industrial Engineering",
	RE3DataSubjectType401:     "Production Technology",
	RE3DataSubjectType40101:   "Metal-Cutting Manufacturing Engineering",
	RE3DataSubjectType40102:   "Primary Shaping and Reshaping Technology",
	RE3DataSubjectType40103:   "Micro-, Precision, Mounting, Joining, Separation Technology",
	RE3DataSubjectType40104:   "Plastics Engineering",
	RE3DataSubjectType40105:   "Production Automation, Factory Operation, Operations Manangement",
	RE3DataSubjectType402:     "Mechanics and Constructive Mechanical Engineering",
	RE3DataSubjectType40201:   "Construction, Machine Elements",
	RE3DataSubjectType40202:   "Mechanics",
	RE3DataSubjectType40203:   "Lightweight Construction, Textile Technology",
	RE3DataSubjectType40204:   "Acoustics",
	RE3DataSubjectType42:      "Thermal Engineering/Process Engineering",
	RE3DataSubjectType403:     "Process Engineering, Technical Chemistry",
	RE3DataSubjectType40301:   "Chemical and Thermal Process Engineering",
	RE3DataSubjectType40302:   "Technical Chemistry",
	RE3DataSubjectType40303:   "Mechanical Process Engineering",
	RE3DataSubjectType40304:   "Biological Process Engineering",
	RE3DataSubjectType404:     "Heat Energy Technology, Thermal Machines, Fluid Mechanics",
	RE3DataSubjectType40401:   "Energy Process Engineering",
	RE3DataSubjectType40402:   "Technical Thermodynamics",
	RE3DataSubjectType40403:   "Fluid Mechanics",
	RE3DataSubjectType40404:   "Hydraulic and Turbo Engines and Piston Engines",
	RE3DataSubjectType43:      "Materials Science and Engineering",
	RE3DataSubjectType405:     "Materials Engineering",
	RE3DataSubjectType40501:   "Metallurgical and Thermal Processes, Thermomechanical Treatment of Materials",
	RE3DataSubjectType40502:   "Sintered Metallic and Ceramic Materials",
	RE3DataSubjectType40503:   "Composite Materials",
	RE3DataSubjectType40504:   "Mechanical Behaviour of Construction Materials",
	RE3DataSubjectType40505:   "Coating and Surface Technology",
	RE3DataSubjectType406:     "Materials Science",
	RE3DataSubjectType40601:   "Thermodynamics and Kinetics of Materials",
	RE3DataSubjectType40602:   "Synthesis and Properties of Functional Materials",
	RE3DataSubjectType40603:   "Microstructural Mechanical Properties of Materials",
	RE3DataSubjectType40604:   "Structuring and Functionalisation",
	RE3DataSubjectType40605:   "Biomaterials",
	RE3DataSubjectType44:      "Computer Science, Electrical and System Engineering",
	RE3DataSubjectType407:     "Systems Engineering",
	RE3DataSubjectType40701:   "Automation, Control Systems, Robotics, Mechatronics",
	RE3DataSubjectType40702:   "Measurement Systems",
	RE3DataSubjectType40703:   "Microsystems",
	RE3DataSubjectType40704:   "Traffic and Transport Systems, Logistics",
	RE3DataSubjectType40705:   "Human Factors, Ergonomics, Human-Machine Systems",
	RE3DataSubjectType408:     "Electrical Engineering",
	RE3DataSubjectType40801:   "Electronic Semiconductors, Components, Circuits, Systems",
	RE3DataSubjectType40802:   "Communication, High-Frequency and Network Technology, Theoretical Electrical Engineering",
	RE3DataSubjectType40803:   "Electrical Energy Generation, Distribution, Application",
	RE3DataSubjectType409:     "Computer Science",
	RE3DataSubjectType40901:   "Theoretical Computer Science",
	RE3DataSubjectType40902:   "Software Technology",
	RE3DataSubjectType40903:   "Operating, Communication and Information Systems",
	RE3DataSubjectType40904:   "Artificial Intelligence, Image and Language Processing",
	RE3DataSubjectType40905:   "Computer Architecture and Embedded Systems",
	RE3DataSubjectType45:      "Construction Engineering and Architecture",
	RE3DataSubjectType410:     "Construction Engineering and Architecture",
	RE3DataSubjectType41001:   "Architecture, Building and Construction History, Sustainable Building Technology, Building Design",
	RE3DataSubjectType41002:   "Urbanism, Spatial Planning, Transportation and Infrastructure Planning, Landscape Planning",
	RE3DataSubjectType41003:   "Construction Material Sciences, Chemistry, Building Physics",
	RE3DataSubjectType41004:   "Structural Engineering, Building Informatics, Construction Operation",
	RE3DataSubjectType41005:   "Applied Mechanics, Statics and Dynamics",
	RE3DataSubjectType41006:   "Geotechnics, Hydraulic Engineering",
}

type RE3DataProviderType string

const (
	RE3DataProviderTypeDataProvider    RE3DataProviderType = "dataProvider"
	RE3DataProviderTypeServiceProvider RE3DataProviderType = "serviceProvider"
)

var RE3DataProviderTypeReverse = map[string]RE3DataProviderType{
	string(RE3DataProviderTypeDataProvider):    RE3DataProviderTypeDataProvider,
	string(RE3DataProviderTypeServiceProvider): RE3DataProviderTypeServiceProvider,
}

type RE3DataAccessType string

const (
	RE3DataAccessTypeOpen       RE3DataAccessType = "open"
	RE3DataAccessTypeEmbargoed  RE3DataAccessType = "embargoed"
	RE3DataAccessTypeRestricted RE3DataAccessType = "restricted"
	RE3DataAccessTypeClosed     RE3DataAccessType = "closed"
)

var RE3DataAccessTypeReverse = map[string]RE3DataAccessType{
	string(RE3DataAccessTypeOpen):       RE3DataAccessTypeOpen,
	string(RE3DataAccessTypeEmbargoed):  RE3DataAccessTypeEmbargoed,
	string(RE3DataAccessTypeRestricted): RE3DataAccessTypeRestricted,
	string(RE3DataAccessTypeClosed):     RE3DataAccessTypeClosed,
}

type RE3DataAccessRestrictions string

const (
	RE3DataAccessRestrictionsFeeRequired             RE3DataAccessRestrictions = "feeRequired"
	RE3DataAccessRestrictionsInstitutionalMembership RE3DataAccessRestrictions = "institutional membership"
	RE3DataAccessRestrictionsRegistration            RE3DataAccessRestrictions = "registration"
	RE3DataAccessRestrictionsFeeOther                RE3DataAccessRestrictions = "other"
)

var RE3DataAccessRestrictionsReverse = map[string]RE3DataAccessRestrictions{
	string(RE3DataAccessRestrictionsFeeRequired):             RE3DataAccessRestrictionsFeeRequired,
	string(RE3DataAccessRestrictionsInstitutionalMembership): RE3DataAccessRestrictionsInstitutionalMembership,
	string(RE3DataAccessRestrictionsRegistration):            RE3DataAccessRestrictionsRegistration,
	string(RE3DataAccessRestrictionsFeeOther):                RE3DataAccessRestrictionsFeeOther,
}

type RE3DataYesNoUn string

const (
	RE3DataYesNoUnYes     RE3DataYesNoUn = "yes"
	RE3DataYesNoUnNo      RE3DataYesNoUn = "no"
	RE3DataYesNoUnUnknown RE3DataYesNoUn = "unknown"
)

var RE3DataYesNoUnReverse = map[string]RE3DataYesNoUn{
	string(RE3DataYesNoUnYes):     RE3DataYesNoUnYes,
	string(RE3DataYesNoUnNo):      RE3DataYesNoUnNo,
	string(RE3DataYesNoUnUnknown): RE3DataYesNoUnUnknown,
}
