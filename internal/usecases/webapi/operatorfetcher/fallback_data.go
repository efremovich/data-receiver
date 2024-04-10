package operatorfetcher

import (
	"github.com/efremovich/data-receiver/internal/entity"
)

// Статичный список операторов на случай, если API для их получения будет недоступно.
var fallbackOperatorList = []entity.Operator{
	{
		Code:       "2AD",
		Name:       "ООО Русь-Телеком",
		IsDisabled: false,
		Thumbs:     []string{"fd04385b9f25ff92bd7e1c5e9eee34d5aad9b3a7"},
	},
	{
		Code:       "2AE",
		Name:       "АО Калуга Астрал",
		IsDisabled: false,
		Thumbs:     []string{"835b6480ce7e40d98f390b53044763b427633806", "3941a1ac98530bc617aaf897840a32df57ee7ff0", "8197aae6bbfb31e831be797a2ca6bde81d794eeb", "2ae075f33e1cea917999ee9271b1923326926acc"},
	},
	{
		Code:       "2AH",
		Name:       "ИнфоТеКС Интернет Траст",
		IsDisabled: false,
		Thumbs:     []string{"91ffb0da8462062c8f2510c369783c433763c0b5", "b6ee40ff320d9b6aa4496b313eab83b58ab96f2b"},
	},
	{
		Code:       "2AK",
		Name:       "ЗАО ТаксНет",
		IsDisabled: false,
		Thumbs:     []string{"62544bdbc5d5098e02641f30464d8223fcb85dc1"},
	},
	{
		Code:       "2AL",
		Name:       "ООО Такском",
		IsDisabled: false,
		Thumbs:     []string{"70c364934655f8344931b02629fc752ba4f1c97b", "3DA77F9618719B9535C85B7CE6D42906D29EAE15"},
	},
	{
		Code:       "2AO",
		Name:       "ООО УЦ АСКОМ",
		IsDisabled: false,
		Thumbs:     nil,
	},
	{
		Code:       "2BA",
		Name:       "АО \"НТЦ СТЭК\"",
		IsDisabled: false,
		Thumbs:     []string{"c5aad32b8fa4c0a490fd32fa1c518e81759a1529"},
	},
	{
		Code:       "2BE",
		Name:       "ООО Компания Тензор",
		IsDisabled: false,
		Thumbs:     []string{"5e78ca1036e036d561fa60bb4616628645eb74f2"},
	},
	{
		Code:       "2BH",
		Name:       "Аргос",
		IsDisabled: true,
		Thumbs:     nil,
	},
	{
		Code:       "2BK",
		Name:       "ООО КОРУС Консалтинг СНГ",
		IsDisabled: false,
		Thumbs:     []string{"1310276925dd08ef0a7b47481d4e512a926b13c0"},
	},
	{
		Code:       "2BM",
		Name:       "ЗАО ПФ СКБ Контур",
		IsDisabled: false,
		Thumbs:     []string{"67aed421936d9f054494d946bf86a88f8a8706db", "5d82568f10ff71241bdf52421c3891f5c23f3834"},
	},
	{
		Code:       "2BN",
		Name:       "ООО Линк-Сервис",
		IsDisabled: false,
		Thumbs:     nil,
	},
	{
		Code:       "2CI",
		Name:       "ООО \"Электронный Экспресс\"",
		IsDisabled: false,
		Thumbs:     []string{"3a543620b8d71d78f2ae750927e45d84d4e104ec", "b4f83eb1ef92cc52472264f166b9d2ccf26c1b31"},
	},
	{
		Code:       "2EE",
		Name:       "ООО «Электронный Экспресс»",
		IsDisabled: false,
		Thumbs:     []string{"74b69d10abc4753b8de3ecde776d5cc683affce6"},
	},
	{
		Code:       "2GS",
		Name:       "ЗАО \"УДОСТОВЕРЯЮЩИЙ ЦЕНТР\"",
		IsDisabled: true,
		Thumbs:     nil,
	},
	{
		Code:       "2HX",
		Name:       "ООО \"Криптэкс\"",
		IsDisabled: false,
		Thumbs:     []string{"c39df25b0e55ca7c6b5dc663bf569b81aebd29d0"},
	},
	{
		Code:       "2IG",
		Name:       "ООО ДИРЕКТУМ",
		IsDisabled: true,
		Thumbs:     []string{"789040b004a089134266c1237ebc5252e2385599", "fe6b6d2e5e5edc8aef11077be85c35e1d2bcae2f"},
	},
	{
		Code:       "2IH",
		Name:       "ООО Э-КОМ",
		IsDisabled: false,
		Thumbs:     []string{"35964809363f265eb6b3643812425b19a7a2b909"},
	},
	{
		Code:       "2IJ",
		Name:       "ООО «Эдисофт»",
		IsDisabled: false,
		Thumbs:     []string{"0820c7ba8c9e1690c8f6d94c6f2921453014c595"},
	},
	{
		Code:       "2IM",
		Name:       "АО \"ЕЭТП\"",
		IsDisabled: false,
		Thumbs:     []string{"0cf33096b56599b4b5f213e2da54c03984d4db2b", "4fa953cc031865688c65eee3a068e70c2242c426"},
	},
	{
		Code:       "2JD",
		Name:       "АО НИИАС",
		IsDisabled: false,
		Thumbs:     []string{"95239eefb1b4dde0216179ad9d57b98716573071", "518dd32f933fd7a18c09528dc73bcb71061ad5a3", "b5ac97f50e72a2b8ffb8ee6531dca8671dbaf8a7"},
	},
	{
		Code:       "2JM",
		Name:       "ООО \"Сислинк\"",
		IsDisabled: false,
		Thumbs:     []string{"8f28d6057fb95d7e80c87baad8803718f395f5a6"},
	},
	{
		Code:       "2KV",
		Name:       "ООО УЦ \"СОЮЗ\"",
		IsDisabled: false,
		Thumbs:     []string{"4ae893adaac55c4bb0303cc26f40d408ed03f688"},
	},
	{
		Code:       "2LB",
		Name:       "ООО ЭТП ГПБ",
		IsDisabled: false,
		Thumbs:     []string{"e12a5fc5330c0bdab79875b19f889d28834899d4"},
	},
	{
		Code:       "2LD",
		Name:       "ООО «ЭЛЕКТРОННЫЕ КОММУНИКАЦИИ»",
		IsDisabled: false,
		Thumbs:     []string{"2f0f972d22fe1b1efb1057fb1ff59e7fb7a75b40"},
	},
	{
		Code:       "2LG",
		Name:       "ООО \"БИФИТ ЭДО\"",
		IsDisabled: false,
		Thumbs:     []string{"daf27231cde82413a960d5b5711e8eb070cbc158"},
	},
	{
		Code:       "2LH",
		Name:       "LERADATA",
		IsDisabled: false,
		Thumbs:     []string{"be604bf64cb463079f97d5ee053ea0217b8ef80c"},
	},
	{
		Code:       "2LJ",
		Name:       "ООО \"Финтендер-крипто\"",
		IsDisabled: false,
		Thumbs:     []string{"61ad8536d58c34df93da3d274e4fa2e0086bcb2b"},
	},
	{
		Code:       "2LT",
		Name:       "ООО \"Оператор-Црпт\"",
		IsDisabled: false,
		Thumbs:     []string{"3fee2553f01a4f55f6564d63959662895ecf287b", "448bb4be6142f8d5800dfc975f011cf151301552", "158b766fb05c249549b5920ab419347864bd1ab6"},
	},
	{
		Code:       "2MA",
		Name:       "Общество с ограниченной ответственностью «ДИАС-К»",
		IsDisabled: true,
		Thumbs:     []string{"5426d27b4f8791159494ff335a131ca899ec542d"},
	},
	{
		Code:       "2MB",
		Name:       "ООО \"Дистэйт\"",
		IsDisabled: false,
		Thumbs:     []string{"f356106e270782e9f5e07aaf25abb4ee8ee9c547", "a6e93ee9c956a96bfe9f2eeea4ac9e56a7886f66", "9f68320fd594764a64080f98f551f5753aa8f708"},
	},
	{
		Code:       "2MD",
		Name:       "ООО «Айтиком»",
		IsDisabled: false,
		Thumbs:     []string{"71302f2a803273b7f32ff9475ebaa95ee0cebc72"},
	},
	{
		Code:       "2PS",
		Name:       "ЭДО.Поток",
		IsDisabled: false,
		Thumbs:     []string{"c791f36ffefa5306c9a6b5aa74531eeb580a2560"},
	},
	{
		Code:       "2VO",
		Name:       "ООО \"Эвотор ОФД\"",
		IsDisabled: false,
		Thumbs:     []string{"4054cb116d742877a1c82633d28d4123e6c223b3"},
	},
	{
		Code:       "BBB",
		Name:       "BBB",
		IsDisabled: false,
		Thumbs:     nil,
	},
	{
		Code:       "XXX",
		Name:       "Контур регрессионного тестирования",
		IsDisabled: false,
		Thumbs:     []string{"2ae075f33e1cea917999ee9271b1923326926acc"},
	},
	{
		Code:       "XXY",
		Name:       "Контур регрессионного тестирования псевдопрямой",
		IsDisabled: false,
		Thumbs:     []string{"2ae075f33e1cea917999ee9271b1923326926acc"},
	},
}
