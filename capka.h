#ifndef CAPKA_H
#define CAPKA_H

int capka_makeKeypair(
	const char*    password, // in
	const char*    saltData, // in
	int            ops,      // in
	int            mem,      // in
	unsigned char* pk,       // out
	unsigned char* sk        // out
);

#endif
