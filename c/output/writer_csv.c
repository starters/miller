#include <stdlib.h>
#include "containers/mixutil.h"
#include "lib/mlrutil.h"
#include "output/writers.h"

typedef struct _writer_csv_state_t {
	int  onr;
	char ors;
	char ofs;
	long long num_header_lines_output;
	slls_t* plast_header_output;
} writer_csv_state_t;

// ----------------------------------------------------------------
// xxx cmt mem-mgmt

static void writer_csv_func(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	writer_csv_state_t* pstate = pvstate;
	char ors = pstate->ors;
	char ofs = pstate->ofs;

	if (pstate->plast_header_output != NULL) {
		// xxx make a fcn to compare these w/o copy: put it in mixutil.
		slls_t* ptemp = mlr_keys_from_record(prec);
		if (!slls_equals(ptemp, pstate->plast_header_output)) {
			slls_free(pstate->plast_header_output);
			pstate->plast_header_output = NULL;
			if (pstate->num_header_lines_output > 0LL)
				fputc(ors, output_stream);
		}
	}

	if (pstate->plast_header_output == NULL) {
		int nf = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			if (nf > 0)
				fputc(ofs, output_stream);
			fputs(pe->key, output_stream);
			nf++;
		}
		fputc(ors, output_stream);
		// xxx maybe mke a mlr_copy_keys_from_record & mlr_reference_keys_from_record
		pstate->plast_header_output = slls_copy(mlr_keys_from_record(prec));
		pstate->num_header_lines_output++;
	}

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputc(ofs, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputc(ors, output_stream);
	pstate->onr++;

	lrec_free(prec); // xxx cmt mem-mgmt
}

static void writer_csv_free(void* pvstate) {
	writer_csv_state_t* pstate = pvstate;
	if (pstate->plast_header_output != NULL) {
		slls_free(pstate->plast_header_output);
		pstate->plast_header_output = NULL;
	}
}

writer_t* writer_csv_alloc(char ors, char ofs) {
	writer_t* pwriter = mlr_malloc_or_die(sizeof(writer_t));

	writer_csv_state_t* pstate = mlr_malloc_or_die(sizeof(writer_csv_state_t));
	pstate->onr                     = 0;
	pstate->ors                     = ors;
	pstate->ofs                     = ofs;
	pstate->num_header_lines_output = 0LL;
	pstate->plast_header_output     = NULL;
	pwriter->pvstate                = (void*)pstate;

	pwriter->pwriter_func = writer_csv_func;
	pwriter->pfree_func   = writer_csv_free;

	return pwriter;
}
