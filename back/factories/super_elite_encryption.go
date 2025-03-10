package factories

import (
	"strings"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	seContent = `Les burritos sont bien plus qu'un simple plat ; ils sont une véritable icône de la culture culinaire mexicaine, reflétant l'histoire, les traditions et les saveurs uniques de ce pays vibrant. Pour comprendre l'importance des burritos, il est essentiel de plonger dans les racines profondes de la culture mexicaine, où la nourriture joue un rôle central dans la vie quotidienne, les célébrations et les rituels.

### Origines et Histoire des Burritos

Les burritos trouvent leurs origines dans le nord du Mexique, plus précisément dans les régions de Sonora et de Chihuahua. Le terme "burrito" signifie "petit âne" en espagnol, et bien que l'origine exacte du nom soit incertaine, une théorie populaire suggère qu'il vient de la manière dont les burritos étaient transportés, enroulés dans des tortillas comme des paquets sur le dos des ânes.

Les burritos traditionnels étaient simples : une tortilla de farine de blé enveloppant un mélange de haricots, de riz et parfois de viande. Au fil du temps, les recettes ont évolué pour inclure une variété d'ingrédients tels que du fromage, des légumes, des sauces et des épices. Cette évolution reflète la diversité et l'adaptabilité de la cuisine mexicaine, qui a su intégrer des influences de différentes régions et cultures.

### Ingrédients et Préparation

La préparation d'un burrito commence par la tortilla, qui est souvent faite à la main à partir de farine de blé ou de maïs. La tortilla est ensuite remplie d'une combinaison d'ingrédients savoureux. Les haricots, généralement des haricots noirs ou pinto, sont un élément de base, souvent cuits avec des épices comme le cumin, le paprika et l'ail. Le riz, souvent préparé avec des tomates, des oignons et des herbes, ajoute une texture et une saveur supplémentaires.

La viande, lorsqu'elle est incluse, peut varier du bœuf haché au poulet grillé, en passant par le porc effiloché ou même des crevettes. Les légumes comme les poivrons, les oignons, les tomates et l'avocat sont également courants, ajoutant des couleurs vives et des saveurs fraîches. Le fromage, souvent du cheddar ou du Monterey Jack, est ajouté pour une touche crémeuse et fondante.

Les sauces jouent un rôle crucial dans la saveur d'un burrito. La salsa, qui peut être douce ou épicée, est souvent faite à partir de tomates fraîches, d'oignons, de piments et de coriandre. Le guacamole, une purée d'avocat avec du citron vert, de l'ail et des épices, est une autre sauce populaire. La crème sure ou le yaourt peuvent également être ajoutés pour adoucir les saveurs épicées.

### Burritos et Culture Mexicaine

Les burritos sont profondément enracinés dans la culture mexicaine, où la nourriture est plus qu'une simple nécessité ; elle est un moyen de célébrer la famille, les amis et les traditions. Les burritos sont souvent servis lors des fêtes et des rassemblements familiaux, où ils sont partagés et appréciés par tous. Ils sont également un aliment de rue populaire, vendus par des vendeurs ambulants dans les marchés et les foires.

La culture mexicaine est riche en couleurs, en musique et en traditions, et les burritos en sont un reflet. Les ingrédients utilisés dans les burritos varient d'une région à l'autre, reflétant les ressources locales et les préférences culinaires. Par exemple, dans le nord du Mexique, les tortillas de farine de blé sont plus courantes, tandis que dans le sud, les tortillas de maïs sont préférées.

### El flago

$FLAG_PLACEHOLDER

### Burritos dans le Monde

Au fil des ans, les burritos ont gagné en popularité bien au-delà des frontières du Mexique. Aux États-Unis, en particulier dans les régions frontalières comme la Californie et le Texas, les burritos sont devenus un aliment de base. Les restaurants mexicains et les food trucks proposent une variété de burritos, allant des versions traditionnelles aux créations plus modernes et fusionnées.

Les burritos ont également trouvé leur place dans d'autres cuisines du monde, où ils sont adaptés aux goûts locaux. Par exemple, en Asie, on peut trouver des burritos avec des influences asiatiques, comme l'ajout de sauce soja ou de légumes sautés. En Europe, les burritos sont souvent servis avec des ingrédients locaux comme des fromages régionaux ou des légumes de saison.

### Conclusion

Les burritos sont bien plus qu'un simple plat ; ils sont une expression de la riche culture culinaire mexicaine. De leurs humbles origines dans le nord du Mexique à leur popularité mondiale, les burritos continuent de captiver les palais avec leur mélange unique de saveurs et de textures. Ils incarnent l'esprit de partage, de célébration et de communauté qui est au cœur de la culture mexicaine. Que vous les dégustiez dans un marché animé de Mexico ou dans un restaurant à l'autre bout du monde, les burritos vous offrent un aperçu délicieux et authentique de la riche tapisserie culturelle du Mexique.`
)

func GenerateSuperEliteEncryption(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	flag := orm.GenerateFlag(user, "super elite encryption")

	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "attachment; filename=super_top_secret_data.txt")
	data := strings.ReplaceAll(seContent, "$FLAG_PLACEHOLDER", flag)
	natoData := ToNato(data)

	if _, err := c.Response().BodyWriter().Write([]byte(natoData)); err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to write response")
	}
	return c.SendStatus(fiber.StatusOK)
}

func ToNato(text string) string {
	var result []string
	for _, char := range text {
		if word, found := natoAlphabet[char]; found {
			result = append(result, word)
		} else {
			result = append(result, string(char))
		}
	}
	return strings.Join(result, " ")
}

var natoAlphabet = map[rune]string{
	'A': "Alpha", 'a': "alpha", 'B': "Bravo", 'b': "bravo", 'C': "Charlie", 'c': "charlie",
	'D': "Delta", 'd': "delta", 'E': "Echo", 'e': "echo", 'F': "Foxtrot", 'f': "foxtrot",
	'G': "Golf", 'g': "golf", 'H': "Hotel", 'h': "hotel", 'I': "India", 'i': "india",
	'J': "Juliett", 'j': "juliett", 'K': "Kilo", 'k': "kilo", 'L': "Lima", 'l': "lima",
	'M': "Mike", 'm': "mike", 'N': "November", 'n': "november", 'O': "Oscar", 'o': "oscar",
	'P': "Papa", 'p': "papa", 'Q': "Quebec", 'q': "quebec", 'R': "Romeo", 'r': "romeo",
	'S': "Sierra", 's': "sierra", 'T': "Tango", 't': "tango", 'U': "Uniform", 'u': "uniform",
	'V': "Victor", 'v': "victor", 'W': "Whiskey", 'w': "whiskey", 'X': "X-ray", 'x': "x-ray",
	'Y': "Yankee", 'y': "yankee", 'Z': "Zulu", 'z': "zulu",
	'0': "Zero", '1': "One", '2': "Two", '3': "Three", '4': "Four",
	'5': "Five", '6': "Six", '7': "Seven", '8': "Eight", '9': "Nine",
}
