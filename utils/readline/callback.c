// +build linux darwin libreadline

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "readline/readline.h"
#include "readline/history.h"

#if !defined(RL_READLINE_VERSION) || (RL_READLINE_VERSION < 0x0500)
typedef char **rl_completion_func_t (const char *, int, int);

#define RL_COMPENTRY_FUNC_T Function
#else
#define RL_COMPENTRY_FUNC_T rl_compentry_func_t
#endif

#include "_cgo_export.h"

char ** call_completion_matches(const char *text) {
    return rl_completion_matches(text, (rl_compentry_func_t *)go_CompletionEntryFunction);
}

void set_completion_entry_function() {
	rl_completion_entry_function = (RL_COMPENTRY_FUNC_T *)go_CompletionEntryFunction;
}

void set_attempted_completion_function() {
	rl_attempted_completion_function = (rl_completion_func_t *)go_AttemptedCompletionFunction;
}

char *null_cstring() {
	return (char *)0;
}

char **null_cstring_array() {
	return (char **)0;
}

char **cstring_array_new(int size) {
	return (char **) malloc(size * sizeof(char *));
}

void cstring_array_set(const char **csa, int i, const char *s) {
	csa[i] = s;
}

const char *cstring_array_get(const char **csa, int i) {
	return csa[i];
}

int cstring_array_len(const char **csa) {
    int n = 0;

    if (csa == (const char **)0)
        return 0;

    while (csa[n] != (char *)0) {
        n++;
    }

    return n;
}
