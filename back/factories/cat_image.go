package factories

import (
	"image"
	"image/color"
	"strings"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/image/bmp"
)

var (
	catImageContent = `The "cat" command in Linux is a versatile and widely-used utility that stands for "concatenate." It is primarily used to read, concatenate, and display the contents of files. Despite its simplicity, "cat" is a powerful tool that can be employed in various scenarios, making it an essential command for both novice and experienced Linux users.

### Basic Usage

The most basic usage of the "cat" command is to display the contents of a file. For example, to view the contents of a file named "example.txt", you would use:

cat example.txt

This command will output the entire contents of "example.txt" to the terminal.

### Concatenating Files

As the name suggests, "cat" can also concatenate multiple files and display their contents sequentially. For instance, to concatenate "file1.txt" and "file2.txt", you would use:

cat file1.txt file2.txt

This will display the contents of "file1.txt" followed by the contents of "file2.txt".

### Redirecting Output

The output of the "cat" command can be redirected to another file using the ">" operator. This is useful for creating new files or overwriting existing ones. For example:

cat file1.txt file2.txt > combined.txt

This command will concatenate "file1.txt" and "file2.txt" and write the result to "combined.txt". If "combined.txt" already exists, it will be overwritten.

To append the contents to an existing file instead of overwriting it, you can use the ">>" operator:

cat file1.txt file2.txt >> combined.txt

### Viewing Line Numbers

The "cat" command can also display line numbers alongside the file contents using the "-n" option:

cat -n example.txt

This will prefix each line of "example.txt" with its corresponding line number.

### Displaying Non-Printable Characters

To view non-printable characters, such as tabs and newlines, you can use the "-v" option:

cat -v example.txt

This will display non-printable characters in a visible format. For example, tabs will be shown as "^I".

### Combining Options

You can combine multiple options to customize the output further. For instance, to display line numbers and non-printable characters, you can use:

cat -nv example.txt

### Using "cat" with Pipes

The "cat" command is often used in conjunction with pipes ("|") to pass the output to other commands. For example, to count the number of lines in a file, you can use:

cat example.txt | wc -l

This will pipe the contents of "example.txt" to the "wc" (word count) command, which will then count and display the number of lines.

### Creating Files

While "cat" is primarily used for reading and concatenating files, it can also be used to create new files. By using the ">" operator with "cat", you can create a new file and enter its contents directly from the terminal. For example:

cat > newfile.txt

After entering this command, you can type the contents of the file. To save and exit, press "Ctrl+D".

### Advanced Usage

#### Displaying Specific Lines

To display specific lines from a file, you can combine "cat" with other commands like "head" or "tail". For example, to display the first 10 lines of a file:

cat example.txt | head -n 10

To display the last 10 lines of a file:

cat example.txt | tail -n 10

#### Searching for Text

You can use "cat" with "grep" to search for specific text within a file. For example, to search for the word "example" in "example.txt":

cat example.txt | grep "example"

This will display all lines containing the word "example."

#### Sorting File Contents

To sort the contents of a file, you can use "cat" with the "sort" command:

cat example.txt | sort

This will sort the lines of "example.txt" alphabetically.

### Practical Examples

#### Example 1: Combining Log Files

Suppose you have multiple log files ("log1.txt", "log2.txt", "log3.txt") and you want to combine them into a single file ("combined_logs.txt"):

cat log1.txt log2.txt log3.txt > combined_logs.txt

#### Example 2: Creating a Backup File

To create a backup of a configuration file ("config.txt") before making changes, you can use:

cat config.txt > config_backup.txt

### Flag

Thank you for reading my blog post about this fantastic command. Here's your flag:

$FLAG_PLACEHOLDER

### Conclusion

The "cat" command is a fundamental tool in the Linux ecosystem, offering a straightforward way to read, concatenate, and manipulate file contents. Its versatility makes it an indispensable command for various tasks, from simple file viewing to complex data processing pipelines. Whether you are a beginner learning the basics of Linux or an experienced user looking to streamline your workflow, mastering the "cat" command can significantly enhance your productivity and efficiency.
`
)

const (
	catImageLongestLine = 483/3 + 1 // used for x, we encode 3 chars per pixel (and the + 1 is for safety)
	catImageLineCount   = 123       // used for y
)

func GenerateCatImage(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	flag := orm.GenerateFlag(user, "cat image")

	cnt := strings.ReplaceAll(catImageContent, "$FLAG_PLACEHOLDER", flag)

	c.Set("Content-Type", "image/bmp")
	c.Set("Content-Disposition", "attachment; filename=cat_image.bmp")
	challengeImage := image.NewRGBA(image.Rect(0, 0, catImageLongestLine, catImageLineCount))

	// Encode the string into the image with char in alpha
	buf := make([]byte, 4)
	for i := 0; i < len(cnt); i++ {
		for j := 0; j < 3; j++ {
			if i+j >= len(cnt) {
				break
			}
			buf[2-j] = cnt[i] >> uint(8*j) & 0xFF
		}
		buf[3] = cnt[i]
		challengeImage.Set(i%catImageLongestLine, i/catImageLongestLine, color.RGBA{buf[0], buf[1], buf[2], buf[3]})
	}

	// Write the image to the response
	err := bmp.Encode(c.Response().BodyWriter(), challengeImage)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to write cat image")
	}

	return c.SendStatus(fiber.StatusOK)
}
