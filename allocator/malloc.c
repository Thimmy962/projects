#define _GNU_SOURCE
#include <unistd.h>
#include <stdio.h>

#ifndef THIMMY
#define THIMMY

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




void *mymalloc(size_t size){
	if(size <= 0)
		return NULL;
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
		heap_end = sbrk(0);
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

int myfree(void *ptr) {
	if((char *)ptr < (char *)heap || (char *)ptr > (char *)heap_end)
		return 1; // an error
	block_t *addr = (block_t *)((char *)ptr - BLOCK_T);
	block_t *pre_block = addr->previous;
	block_t *next_block = addr->next;
	if(pre_block != NULL && pre_block->free == TRUE) {

		pre_block->next = next_block;
		pre_block->size = pre_block->size + BLOCK_T + addr->size;
		if(next_block != NULL)
			next_block->previous = pre_block;

		addr = pre_block;
	}

	if(next_block != NULL && next_block->free == TRUE) {
		block_t *another_next = next_block->next;
		addr->next = another_next;
		if(another_next != NULL)
			another_next->previous = addr;
		addr->size = addr->size + BLOCK_T + next_block->size;
		addr->free = TRUE;
	}
	else {
		addr->free = TRUE;
	}
	if((char *)next_free > (char *)addr)
		next_free = addr;
	
	return 0; //success
}

#endif

