#include <sys/select.h>


unsigned long int __fdelt_chk (unsigned long int d)
{
	if (d >= FD_SETSIZE)
		__chk_fail ();

	return d / __NFDBITS;
}

