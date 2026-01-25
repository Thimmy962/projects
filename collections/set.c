/**
 * implementing a set for ints using open addressing of a hashtable
 */
#include <stdint.h>
#include <stdlib.h>

 typedef struct{
    int **array;
    int16_t size;
    int16_t capacity;
}Set;

void destroy_set(Set *);
int16_t hash(Set *, int);
void add(Set *, int);
Set *init_set();
void resize(Set *);
void printall(Set *);
void delete(Set *, int);

Set *init_set() {
    Set *set = malloc(sizeof(Set));
    if(!set) return NULL;
    set->capacity = 2;
    set->size = 0;

    set->array = malloc((size_t)sizeof(*set->array) * (size_t)set->capacity);
    if(!set->array) {
        free(set);
        return NULL;
    }

    for(int i = 0; i < set->capacity; i++)
        set->array[i] = NULL;
    
    return set;
}

int16_t hash(Set *set, int data) {
    return (int16_t)data % set->capacity;
}

void add(Set *set, int value) {
    for(int i = 0; i < set->capacity; i++) {
        if(set->array[i] != NULL && *set->array[i] == value)
            return;
    }
    int *data = malloc(sizeof(value));
    
    if(!data)
    {
        destroy_set(set);
        exit(1);
    }
    *data = value;

    int index = hash(set, value);

    if(!set->array[index]) set->array[index] = data;
    else {
        int iter = ++index;
        while(iter != index) {
            if(iter == set->capacity) iter = iter % set->capacity;
            if(!set->array[iter]) {
                set->array[iter] = data;
                break;
            }
            iter++;
        }
    }
       
    ++set->size;

    if((float)set->size/(set->capacity) > 0.4)
        resize(set);

}

void destroy_set(Set *set) {
    for(int i = 0; i < set->capacity; i++) {
        if(set->array[i])
            free(set->array[i]);
    }
    free(set->array);
    free(set);
}


void resize(Set *set) {
    int32_t current_capcity = set->capacity;
    set->capacity *= 2;
    int **new_mem = malloc((size_t)sizeof(set->array[0]) * (size_t)set->capacity);
    if(!new_mem) {
        destroy_set(set);
        exit(1);
    }
    int i = 0;
    for(; i < current_capcity; i++) {
        new_mem[i] = set->array[i];
    }
    /** since i is equal to current_capacity
        if the size of the old memory is 5, size of new memory is 10
        copy everything from 0-4 of old to 0-4 of new;
        init 5-9 of new to null
    */
    for(; i < set->capacity; i++) {
        new_mem[i] = NULL;
    }

    int **tmp = new_mem;
    new_mem = set->array;
    set->array = tmp;

    free(new_mem);
}


void delete(Set *set, int value) {
    for(int i = 0; i < set->capacity; i++) {
        if(set->array[i] && *set->array[i] == value){
            int *val = set->array[i];
            free(val);
            set->array[i] = NULL;
            set->size--;
            return;
        }
    }
}

void printall(Set *set) {
    for(int i = 0; i < set->capacity; i++) {
        if(set->array[i]) {
            printf("%d\n", *set->array[i]);
        }
    }
}