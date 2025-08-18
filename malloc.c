#define _GNU_SOURCE
#include <unistd.h>
#include <stdio.h>

#define TRUE 1
#define FALSE 0

typedef struct meta_data {
	size_t size;
	struct meta_data *next;
	struct meta_data *previous;
	size_t free;
}block_t;


#define PAGE 4096 
#define BLOCK_T sizeof(block_t)



block_t *heap = NULL; // beginning of the heap
block_t *heap_end = NULL; // end of the head
block_t *next_free = NULL; // free address closest to the heap




void *malloc(size_t size){
	if(heap == NULL){
		heap = sbrk(PAGE);
		heap->size = PAGE - BLOCK_T;
		heap->next=heap->previous=NULL;
		heap->free=TRUE;
		heap_end = sbrk(0);
		next_free = heap;
	}

	size_t total = BLOCK_T + size;
	block_t *new = next_free;
	while(new) {
    	if(new->free && new->size >= total)
        	break;  // found a suitable block
    	new = new->next;
	}
	
	if(new && new->next == NULL && new->free && new->size < total) {
    	size_t more = PAGE; // new is the last block (tail), but too small
    	sbrk(more);
    	new->size += more;   // extend its size
	}

	// create a  block * for the new fragment and init its value from the old block  	
	block_t *next = (block_t *)((char *)new + total);
	next->free=TRUE;
	next->size = new->size - total;
	next->next = new->next;
	next->previous = new;

	new->next = next;
	new->size = size;
	new->free = FALSE;

	if((char *)next_free > (char *)next)
		next_free = next;

	return (char *)new + BLOCK_T;	
}
