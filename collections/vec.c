#ifndef VECTOR
#define VECTOR
/**
 * This is my implementation of dynamic arrays for an int in C as a stack
 */

#include <stdint.h>
#include <stdlib.h>
#include <string.h> 
/**
 * vecerr is used to define errors in this lib
 * always check that vecerr is 0 before proceeding else there is an error
 * 
*/ 


 int vecerr;

typedef struct  {
    int *array;
    int16_t size;
    int16_t capacity;
    int16_t init_capacity;
} Vector;

typedef struct VectorOps {
    void (*vec_push)(Vector *, int16_t);
    int  (*vec_get)(Vector *, int16_t);
} VectorOps;

void expand(Vector *v) {
    v->capacity *= 2;
    int *array = v->array;
    vecerr = 0;
    v->array = malloc(v->capacity * sizeof(int));
    if(v->array == NULL) {
        vecerr = 1;
       return; 
        
    }

    memcpy(v->array, array, v->size);
    memset(array, 0, v->size);
    free(array);
    array = NULL;

}

Vector *init_vect(int16_t size) {
    Vector *vec = malloc(sizeof(Vector));
    vecerr = 0;
    if(size == 0)
        size = 5;
    if(vec == NULL) {
        vecerr = 1;
        return NULL;
    }
    vecerr = 0;
    int *intaddr = malloc(sizeof(int) * size);
    if(intaddr == NULL) {
        vecerr = 1;
        return NULL;
    }
    memset(intaddr, 0, size);
    vec->array = intaddr;
    vec->init_capacity = vec->capacity = size;
    vec->size = 0;

    return vec;
}

void shrink(Vector *vec) {
    vec->capacity /= 2;
    vecerr = 0;
    int *array = malloc(vec->capacity);
    if(array == NULL) {
        vecerr = 1;
        return;
    }
    memcpy(vec->array, array, vec->size);
    memset(array, 0, vec->size);
    free(array);
    array = NULL;
}


void push(Vector *v, int16_t value) {
    if(v->capacity == v->size) {
        expand(v);
        if(vecerr != 0)
            return;
    }
    if(v->size < v->capacity/2 && v->size > v->init_capacity) {
        shrink(v);
        if(vecerr != 0)
            return;
    }

    v->array[v->size++] = value;
}

int get(Vector *vec, int16_t position) {
    vecerr = 0;
    if(position < 0 || position >= vec->size) {
        vecerr = -1;
        return -1000000;
    }
    return vec->array[position];
}


VectorOps methods = {
    .vec_get = get,
    .vec_push = push
};
#endif