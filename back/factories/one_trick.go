package factories

import (
	"bytes"
	"math/rand/v2"
	"strings"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	otChalCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	otChalContent = `## **Journal de Bord du Capitaine Icarus L. Vortex**  
**Date : 17e cycle lunaire de l’an 4021**  
**Coordonnées : 47°72’X - 18°92’Y dans le secteur d’Andromède-Zeta**  

---

### **Entrée 1 : L’Odyssée du Vide**  
Il y a maintenant 127 jours que nous avons quitté la station mère. Le vaisseau *Oblivion* continue sa route dans l'immensité du néant cosmique. La gravité artificielle tient bon, mais les derniers ajustements ont révélé une anomalie dans la chambre cryogénique. L'IA du bord, S.A.P.H.I.R.A., m'a assuré que ce n’était qu’un léger dérèglement... mais je n’en suis pas convaincu.  

L'équipage commence à ressentir les effets de l'isolement. L'ingénieur Vega passe de plus en plus de temps enfermé dans la salle des machines, marmonnant des calculs incompréhensibles. Le docteur Nyméria a quant à elle entrepris de cataloguer chaque grain de poussière qui se dépose sur les hublots...  

Le journal de bord se doit de contenir tous les éléments nécessaires pour la postérité. Qui sait, peut-être que dans mille ans, une civilisation extraterrestre tombera sur ces enregistrements ?  

---

### **Entrée 12 : L’Enigme du CryoBloc-7**  
S.A.P.H.I.R.A. a relevé une signature énergétique anormale provenant du compartiment CryoBloc-7. Nous avons procédé à un scan approfondi, mais les résultats sont illisibles. Seuls des fragments de code crypté apparaissent à l’écran :  

> U0dWeWMybHVaU0JrYVhJdVkyOXRMMkZ3Y0hKdmRHTXRaVzUwWlc1MEwyTmhjM05wYjI0dWIzSm5MMkpwYjI0PQ==  

J’ai tenté de décrypter cette séquence, mais en vain. Peut-être un vestige des protocoles de sécurité embarqués...  

---

### **Entrée 27 : Une étrange découverte**  
Hier soir, alors que nous traversions la ceinture d’astéroïdes de Kepler-442, Vega a signalé une pulsation électromagnétique d'origine inconnue. Cela ressemblait à un signal de détresse... ou un avertissement.  

Nous avons récupéré une boîte métallique flottant dans l’espace, marquée d’un symbole que je ne reconnais pas. À l’intérieur, un simple fragment de papier où étaient griffonnées ces lettres :  

**$FLAG_PLACEHOLDER**

J’ignore encore ce que cela signifie. Peut-être n’est-ce qu’un leurre, un message perdu depuis des siècles. Mais une chose est sûre... quelqu’un, quelque part, voulait que ce message soit trouvé.  
`
)

// Returns the user's XOR key for the challenge
func otKeyGenerator(user *orm.User) []byte {
	var key bytes.Buffer
	rndSrc := rand.NewPCG(user.RandomSeed, orm.AsciiSum("one trick key"))
	for i := 0; i < 32; i++ {
		key.WriteByte(otChalCharset[rndSrc.Uint64()%uint64(len(otChalCharset))])
	}
	return key.Bytes()
}

func RenderOneTrick(c *fiber.Ctx) error {
	return c.Render("exclusive_club/index", fiber.Map{})
}

func SubmitOneTrick(c *fiber.Ctx) error {
	key := c.FormValue("key")

	if len(key) != 32 {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "Invalid key length")
	}

	user := middlewares.ReadUser(c)
	userKey := otKeyGenerator(user)
	flag := orm.GenerateFlag(user, "exclusive club")
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", "filename=decrypted_secret_journal.txt")

	userContent := strings.ReplaceAll(otChalContent, "$FLAG_PLACEHOLDER", flag)
	userContentBytes := []byte(userContent)
	var cipherBuf bytes.Buffer
	cipherBuf.Grow(len(userContentBytes))
	for i, b := range userContentBytes {
		cipherBuf.WriteByte(b ^ userKey[i%len(userKey)] ^ key[i%len(key)])
	}

	c.WriteString(cipherBuf.String())
	return c.SendStatus(fiber.StatusOK)
}

func EncryptOneTrick(c *fiber.Ctx) error {
	text := c.FormValue("text")
	if len(text) > 256 {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "Text too long")
	}
	user := middlewares.ReadUser(c)
	userKey := otKeyGenerator(user)
	textBytes := []byte(text)
	var cipherBuf bytes.Buffer
	cipherBuf.Grow(len(textBytes))
	for i, b := range textBytes {
		cipherBuf.WriteByte(b ^ userKey[i%len(userKey)])
	}
	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "filename=encrypted_text.txt")
	c.WriteString(cipherBuf.String())
	return c.SendStatus(fiber.StatusOK)
}
