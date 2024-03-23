#include <sodium.h>
#include <string.h>

int capka_makeKeypair(
	const char*    password, // in
	const char*    saltData, // in
	int            ops,      // in
	int            mem,      // in
	unsigned char* pk,       // out
	unsigned char* sk        // out
) {
	unsigned char salt[crypto_pwhash_SALTBYTES] = {0};
	if(crypto_generichash(
		   salt, sizeof(salt), (unsigned char*) saltData, strlen(saltData),
		   NULL, 0
	   ) != 0)
		return -1;

	unsigned char seed[crypto_sign_SEEDBYTES] = {0};
	if(crypto_pwhash(
		   seed, sizeof(seed), password, strlen(password), salt, ops, mem,
		   crypto_pwhash_ALG_ARGON2ID13
	   ) != 0)
		return -1;

	if(crypto_sign_seed_keypair(pk, sk, seed) != 0) return -1;

	return 0;
}
