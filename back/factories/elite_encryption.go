package factories

import (
	"encoding/base64"
	"strings"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	chalContent = `Modern cryptography is a vast and complex field that has evolved significantly over the years, driven by advances in technology, computer science, and mathematics. At its core, cryptography is the practice of using mathematical algorithms to secure data, communications, and transactions. This involves transforming plaintext (readable data) into unreadable ciphertext (unreadable data) using a secret key or password.

The history of modern cryptography dates back to World War II, when the British government established the Government Code and Cypher School (GC&CS), now known as GCHQ (Government Communications Headquarters). The GC&CS was tasked with breaking enemy codes and ciphers, which laid the foundation for the development of cryptography itself. In response, cryptographers began to develop their own cryptographic systems that could remain secret from potential adversaries.

In the 1970s, public-key cryptography emerged as a major breakthrough in cryptography. This type of cryptography uses pairs of keys: a public key and a private key. The public key is shared openly with others, while the private key remains secret. Any data encrypted using the public key can be decrypted only by someone possessing the corresponding private key. Public-key cryptography revolutionized the way cryptography was used, enabling secure communication over the internet.

One of the most influential cryptographic algorithms in modern times is RSA (Rivest-Shamir-Adleman). Developed in 1978, RSA is a public-key encryption algorithm that uses the difficulty of factoring large numbers to ensure security. The basic idea behind RSA is as follows:

1. Key generation: Two large prime numbers are generated, which are then multiplied together to produce a composite number (n).
2. Public key: n and a related value e are made available publicly.
3. Private key: The private key d is computed such that d*e = 1.

The security of RSA relies on the difficulty of factoring large composite numbers into their prime factors. An attacker attempting to factor n will need to try all possible combinations, which makes it computationally infeasible to break the encryption.

Other important cryptographic algorithms include:

* AES (Advanced Encryption Standard): a symmetric-key block cipher that is widely used for encrypting data at rest and in transit.
* SHA-256 (Secure Hash Algorithm 256): a cryptographic hash function that produces a fixed-size, unique output for any input message.
* Diffie-Hellman key exchange: a method for securely exchanging secret keys over an insecure channel.

Modern cryptography also relies heavily on quantum computing-resistant algorithms. As quantum computers become more powerful, they will be able to break many classical encryption algorithms in a matter of seconds. To mitigate this threat, researchers have developed cryptographic algorithms that are resistant to quantum attacks.

One such algorithm is lattice-based cryptography, which uses the hardness of problems related to lattice structures (e.g., shortest vectors) to ensure security. Lattice-based cryptosystems are thought to be resistant to quantum attacks, as they do not rely on large prime factorization or discrete logarithms.

Another important area in modern cryptography is cryptographic protocols and standards. The Internet Engineering Task Force (IETF) has developed several key protocols for secure communication over the internet, including:

* TLS (Transport Layer Security): a cryptographic protocol that provides secure communication between web browsers and servers.
* SSH (Secure Shell): a cryptographic protocol that provides secure access to remote computers.
* IPsec (Internet Protocol Security): a suite of cryptographic protocols that provide secure communication over the internet.

The IETF also has developed several key standards for cryptography, including:

* NIST SP 800-56A: a standard for key management and encryption.
* FIPS 140-3: a standard for cryptographic modules and devices.
* RFC 6080: a standard for AES-GCM (Galois/Counter Mode).

Modern cryptography also relies on advanced mathematical concepts, such as:

* Elliptic Curve Cryptography (ECC): an algorithm that uses the properties of elliptic curves to provide secure encryption and digital signatures.
* Homomorphic Encryption: an algorithm that allows computations to be performed on encrypted data without decrypting it first.

The future of modern cryptography is likely to involve a number of exciting developments, including:

* Quantum-resistant algorithms: cryptographic algorithms designed to resist attacks by quantum computers.
* Post-quantum cryptography: the study of cryptographic algorithms that are resistant to quantum attacks and have not yet been broken.
* Side-channel attack resistance: the development of cryptographic algorithms that can withstand side-channel attacks (e.g., timing, power consumption).
* Secure multi-party computation: the ability to perform computations on private data in a secure and private manner.

In conclusion, modern cryptography is a rapidly evolving field that relies on advances in mathematics, computer science, and technology. From public-key cryptography to quantum-resistant algorithms, modern cryptography has become an essential component of online security. As technology continues to advance, we can expect to see new cryptographic algorithms and protocols emerge that will further secure our digital lives.

In terms of specific applications of modern cryptography, there are many examples:

* Online banking: encryption is used to protect financial data and ensure secure transactions.
* E-commerce: encryption is used to protect customer data and ensure secure transactions.
* Healthcare: encryption is used to protect sensitive patient data and ensure secure communication between healthcare providers.
* Secure messaging apps: encryption is used to protect user communications and ensure secure authentication.

In terms of real-world threats, modern cryptography must contend with a variety of attacks and challenges. Some examples include:

* Quantum computing attacks: the ability of quantum computers to break certain classical encryption algorithms.
* Side-channel attacks: the use of timing, power consumption, or other side effects to compromise cryptographic security.
* Insider threats: the intentional or unintentional disclosure of sensitive information by authorized personnel.

To mitigate these threats, researchers and developers are working on a number of solutions, including:

* Quantum-resistant algorithms: cryptographic algorithms designed to resist attacks by quantum computers.
* Side-channel attack resistance: the development of cryptographic algorithms that can withstand side-channel attacks.
* Secure multi-party computation: the ability to perform computations on private data in a secure and private manner.

In conclusion, modern cryptography is a rapidly evolving field that has become an essential component of online security. From public-key cryptography to quantum-resistant algorithms, modern cryptography continues to advance our understanding of secure communication and data protection.

RSA (Rivest-Shamir-Adleman) is a widely used public-key encryption algorithm that has been the cornerstone of secure communication over the internet for decades. Here's a detailed overview of how RSA works and its current status:

**History**

RSA was first developed in 1978 by Ron Rivest, Adi Shamir, and Leonard Adleman at MIT. The algorithm was initially designed to be unbreakable using known computer algorithms at the time. Since then, RSA has become one of the most widely used encryption algorithms in the world.

**Key Generation**

The key generation process for RSA involves generating two large prime numbers, p and q, which are then multiplied together to produce a composite number n = p*q. The value of n is the modulus used for encrypting and decrypting data.

To generate the public key, an algorithm called the Euler's totient function φ(n) is used. φ(n) is calculated as (p-1)*(q-1), which gives the total number of possible values that can be used to encrypt a message.

The public key consists of two components:

* n: the modulus
* e: the public exponent, which is an integer less than φ(n)

The private key consists of two components:

* n: the modulus (same as in the public key)
* d: the private exponent, which is the modular inverse of e modulo φ(n)

**Encryption**

To encrypt a message using RSA, the following steps are taken:

1. Convert the plaintext message to numerical values.
2. Calculate the ciphertext values by raising each numerical value to the power of e modulo n.

The resulting ciphertext can be decrypted using the corresponding private key (d).

**Decryption**

To decrypt the ciphertext using RSA, the following steps are taken:

1. Raise each ciphertext value to the power of d modulo n.
2. Convert the resulting numerical values back to plaintext.

**Security**

RSA's security is based on the difficulty of factoring large composite numbers into their prime factors (p and q). An attacker attempting to factor n will need to try all possible combinations, which makes it computationally infeasible to break the encryption.

However, RSA has some limitations:

* Key size: The larger the key size, the more secure the algorithm is. However, increasing key size also increases computational overhead.
* Factorization attacks: Attackers can use specialized algorithms and hardware to factor large numbers, which can compromise RSA's security.

**Current Status**

RSA remains widely used due to its simplicity and broad support across various platforms. However, as quantum computing becomes more prevalent, RSA's security is increasingly threatened.

In 2017, Google announced that they had broken RSA keys of size 2048 bits using their " Quantum Approximate Optimization Algorithm" (QAOA). This was a significant milestone, as it demonstrated the feasibility of breaking RSA with a small-scale quantum computer.

To address this threat, researchers have developed new algorithms and techniques to improve RSA's security, such as:

* Key stretching: Adding extra computational steps to slow down factorization attacks.
* Hybrid cryptosystems: Combining RSA with other encryption algorithms to provide additional security.
* Quantum-resistant key sizes: Developing key sizes that are resistant to quantum attacks.

In conclusion, RSA is a widely used public-key encryption algorithm with a rich history. While it remains secure for many applications, its vulnerability to factorization attacks and increasing threat from quantum computing necessitates ongoing research into improving its security.

**Real-World Applications**

RSA has numerous real-world applications:

* Online banking: RSA is often used to encrypt sensitive data transmitted between the bank's servers and users' browsers.
* Secure messaging apps: RSA is used for end-to-end encryption in popular messaging services like WhatsApp, Signal, and Telegram.
* Virtual private networks (VPNs): RSA can be used to secure VPN connections by encrypting data exchanged between users.

**Alternatives**

There are several alternatives to RSA that offer improved security:

* Elliptic Curve Cryptography (ECC): ECC is a more efficient algorithm with smaller key sizes while maintaining equivalent security.
* Diffie-Hellman key exchange: This algorithm provides secure key exchange without relying on public-key cryptography.
* Quantum-resistant algorithms: New algorithms like lattice-based cryptography and code-based cryptography are being developed to provide quantum resistance.

In summary, RSA remains a widely used encryption algorithm due to its simplicity and broad support. However, as quantum computing advances, it is essential to explore alternative solutions that can provide improved security for sensitive data transmission.

Hashing algorithms are a crucial component of modern cryptography, allowing us to verify the integrity and authenticity of digital data. In this section, we'll explore the basics of hashing, the limitations of MD5, and introduce some modern alternatives like SHA-256 and SHA-512.

**What is Hashing?**

Hashing is a one-way process that transforms large amounts of data into a fixed-size string of characters, known as a hash value or digest. This process is irreversible, meaning it's impossible to obtain the original data from the hash value. Hashing is used for various purposes, including:

* Data integrity verification: ensuring that data hasn't been tampered with during transmission or storage.
* Authentication: verifying the identity of a digital entity, such as a user or device.
* Non-repudiation: proving that a sender has sent a message without being able to deny it.

**MD5 (Message-Digest Algorithm 5)**

MD5 was one of the first widely used hashing algorithms, introduced in 1991. It's a cryptographic hash function with a fixed output size of 128 bits. MD5 is fast and efficient but has several significant limitations:

* **Collision vulnerability**: MD5 can produce different outputs for the same input (a "collision").
* **Preimage attacks**: It's possible to find an input that produces a specific output (a "preimage").

In 2008, the MD5 collision was discovered, which means it's now considered deprecated. While MD5 is still in use in some legacy systems, it's no longer recommended for new applications due to its security vulnerabilities.

**SHA-256 and SHA-512**

The Secure Hash Algorithm (SHA) family is a widely used set of cryptographic hash functions developed by the National Security Agency (NSA). The most commonly used algorithms are:

* **SHA-256**: A 256-bit hash function with a fast computation time. It's considered secure for most purposes, including data integrity verification and digital signatures.
* **SHA-512**: A 512-bit hash function with an even slower computation time than SHA-256. However, it provides additional security benefits due to its increased output size.

**Other Modern Hashing Algorithms**

Some other modern hashing algorithms worth mentioning are:

* **BLAKE2b**: A high-speed hash function designed for cryptographic applications.
* **Keccak-256**: A widely used hash function in blockchain protocols like Ethereum and IPFS.
* **RIPEMD-160**: A hash function developed by the RIPEAL project.

**Comparison of Hashing Algorithms**

Here's a brief comparison of some popular hashing algorithms:

| Algorithm | Output size (bits) | Speed (hash per second) |
| --- | --- | --- |
| MD5 | 128 | Fast (~10^2 hashes/s) |
| SHA-256 | 256 | Moderate (~10^3 hashes/s) |
| SHA-512 | 512 | Slow (~10^4 hashes/s) |
| BLAKE2b | Varies (up to 384 bits) | Very fast (~10^7 hashes/s) |

When choosing a hashing algorithm, consider the following factors:

* **Security**: SHA-256 and SHA-512 are considered secure for most purposes.
* **Speed**: BLAKE2b and SHA-256 offer good balance between speed and security.
* **Output size**: Larger output sizes like SHA-512 provide additional security benefits.

**Best Practices**

To ensure the integrity and authenticity of your data, follow these best practices:

1. Use a secure hashing algorithm like SHA-256 or SHA-512.
2. Use a random salt value to prevent collisions.
3. Store hash values securely, away from plaintext data.
4. Regularly update and rotate your secret keys.

In conclusion, hashing algorithms play a critical role in maintaining the integrity and authenticity of digital data. While MD5 has its limitations, modern alternatives like SHA-256 and SHA-512 offer improved security benefits. By choosing the right algorithm for your use case and following best practices, you can ensure the secure storage and transmission of sensitive information.

**PBKDF2 (Password-Based Key Derivation Function 2)**

PBKDF2 is a widely used password-based key derivation function that was designed to provide slow and computationally expensive hashing, making it resistant to brute-force attacks. It's often used in conjunction with other algorithms like AES for encrypting data.

**Key Features**

* **Slow hashing**: PBKDF2 uses a slow hashing algorithm to make it computationally expensive to calculate the hash value.
* **Key stretching**: The algorithm stretches the password into a longer key, making it harder to guess or crack.
* **Salted**: PBKDF2 uses a random salt value to prevent rainbow table attacks.

**PBKDF2 Variants**

There are several variants of PBKDF2, including:

* $FLAG_PLACEHOLDER
* **PBKDF2-SHA-1**: Uses SHA-1 hashing algorithm (depreciated due to its security vulnerabilities)
* **PBKDF2-HMAC-SHA-256**: Uses HMAC with SHA-256 hashing
* **PBKDF2-HMAC-SHA-512**: Uses HMAC with SHA-512 hashing

**Argon2**

Argon2 is a password-hashing algorithm designed by Google and introduced in 2015. It's considered one of the most secure password-hashing algorithms available.

**Key Features**

* **Memory-hard**: Argon2 uses memory-hard functions, which make it resistant to GPU-based attacks.
* **Time-slicing**: The algorithm is time-sliced, which means that the hash value is calculated over a fixed amount of time, making it harder to exploit with specialized hardware.
* **Randomized computational cost**: Argon2's computational cost varies depending on the system's memory and CPU resources.

**Comparison of PBKDF2 and Argon2**

Here's a brief comparison of PBKDF2 and Argon2:

| Algorithm | Output size (bits) | Computational overhead |
| --- | --- | --- |
| PBKDF2-HMAC-SHA-256 | 256 | Medium (~10^3 operations/s) |
| PBKDF2-HMAC-SHA-512 | 512 | High (~10^4 operations/s) |
| Argon2 | Varies (up to 1024 bits) | Very high (~10^6 operations/s) |

**Best Practices**

When choosing a password-hashing algorithm, consider the following factors:

* **Memory constraints**: If memory is limited, Argon2's memory-hard design may be beneficial.
* **CPU constraints**: PBKDF2 and Argon2 can both be used on CPUs with moderate computational power.
* **Security requirements**: Argon2 is considered one of the most secure password-hashing algorithms available.

To use these algorithms securely:

1. **Use a slow hashing algorithm**: Choose an algorithm that provides sufficient computational overhead to prevent brute-force attacks.
2. **Store salt values securely**: Store salt values separately from plaintext data and ensure they're not easily guessable.
3. **Regularly update your password storage**: Rotate passwords regularly and use the chosen algorithm to hash new passwords.

In conclusion, PBKDF2 and Argon2 are both widely used password-hashing algorithms with different strengths and weaknesses. When choosing an algorithm, consider factors like memory constraints, CPU constraints, and security requirements. By following best practices and using these algorithms securely, you can protect sensitive information from unauthorized access.`
)

func GenerateEliteEncryption(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	flag := orm.GenerateFlag(user, "elite encryption")

	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "attachment; filename=top_secret_data.txt")
	data := strings.ReplaceAll(chalContent, "$FLAG_PLACEHOLDER", flag)
	writer := base64.NewEncoder(base64.StdEncoding, c.Response().BodyWriter())
	defer writer.Close()
	if _, err := writer.Write([]byte(data)); err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to write response")
	}
	return c.SendStatus(fiber.StatusOK)
}
