#ifndef HASHTABLE
#define HASHTABLE

/**
 * implementing a hashmap using separate chaining for an int.
 * When any bucket has 10 children, the hashmap increases in size and rehash all the structs in all the buckets
 */

#include <stdint.h>
#include <stdlib.h>
#include <string.h>
/**
* mapErr is the error value in this map equivalent to errno in global C
*/
int mapErr;

typedef struct Node{
	char *key;
	void *value;
	struct Node *next;
} Node;


typedef struct {
	Node **array;
	int32_t size;
	int32_t filled;
}Map;


Map *init_map();
void insert(Map *, char *, void *);
int hash_function(Map *, char *);
void resize(Map *);
void *get(Map *, char *);
void add(Map *map, char *key, void *value);
int hash(Map *map, char *key);

typedef struct {
	Map *(*init_map)();
	void (*insert)(Map *, char *, void *);
	int (*hash_function)(Map *, char *);
	void (*resize)(Map *);
	void *(*get)(Map *, char *);
}MapOps;

MapOps methods = {
	.init_map = init_map,
	.insert = add,
	.hash_function = hash,
	.get = get
	
};

Map *init_map() {
	Map *map = malloc(sizeof(Map));
	mapErr = 0;
	if(map == NULL) {
		mapErr = -1;
		perror("Init error");
		exit(1);
	}
	map->size = 50;
	/**
	 * an array where each index stores the pointer to Node *
	 * [Node *|Node *|Node *|Node *|Node *|Node *|Node *|Node *|Node *]
	 * 100 contigous 
	 */
	int32_t size = map->size * (int32_t)sizeof(*map->array);
	map->array = malloc((size_t)size);
	if(!map->array) {
		perror("Init error");
		exit(1);
	}
	map->filled = 0;

	for(int i = 0; i < map->size; i++) {
		map->array[i] = NULL;
	}
	
	return map;
}

void add(Map *map, char *key, void *value) {
	int index = methods.hash_function(map, key);
	Node *node = malloc(sizeof(Node));
	if(node == NULL)
		exit(1);
	node->key = key;
	node->next = NULL;
	node->value = value;


	Node *address = map->array[index];
	Node *parent = NULL;
	for(Node *tmp = address; tmp != NULL; tmp = tmp->next)
		parent = tmp;

	if(parent == NULL)
		map->array[index] = node;
	else
		parent->next = node;
	if((float)++map->filled / (float)map->size > 0.5) {
		resize(map);
	}
}

void addNode(Map *map, Node *node) {
	int index = methods.hash_function(map, node->key);
	Node *addr = map->array[index], *parent = NULL;
	for(Node *tmp = addr; tmp != NULL; tmp = tmp->next)
		parent = tmp;
	if(!parent)
		map->array[index] = node;
	else
		parent->next = node;
	
	if((float)++map->filled / (float)map->size > 0.75)
		resize(map);
}

/**
 * simple hash function for the hashmap
 */
int hash(Map *map, char *key) {
	mapErr = 0;
	if(key == NULL) {
		mapErr = -1;
	}
	
	int val = 0;
	char *s = key;
	while(*s != '\0') {
		val += (*s - '0');
		s++;
	}
	return val % map->size;
}


void resize(Map *map) {
	int32_t current_size = map->size;
	Node **tmp = map->array;
	map->size = map->size * 2;
	// request a new memory for the list (double the size of the previous list)
	map->array = malloc((size_t)map->size * sizeof(Node *));
	if(!map->array) {
		perror("Could not allocate memory for resizing");
		exit(1);
	}
	// set the indices of the new list to NULL
	for(int i = 0; i < map->size; i++) {
		map->array[i] = NULL;
	}
	for(int i = 0; i < current_size; i++) {
		if(tmp[i] != NULL) {
			Node *node = tmp[i];
			Node *tmpp = node;
			while(tmpp->next != NULL) {
				node = tmpp;
				tmpp = tmpp->next;
				node->next = NULL;
				addNode(map, node);
			}
			addNode(map, tmpp);
		}
	}
}

void *get(Map *map, char *key) {
	int32_t index = methods.hash_function(map, key);
	
	Node *tmp = map->array[index];
	for(; tmp != NULL; tmp = tmp->next) {
		if(strcmp(tmp->key, key) == 0)
			return tmp->value;
	}
	return NULL;
}

void del(Map *map, char *key) {
	int32_t index = methods.hash_function(map, key);
	
	Node *tmp = map->array[index], *parent = NULL;
	for(; tmp != NULL; tmp = tmp->next) {
		if(strcmp(tmp->key, key) == 0) {
			if(parent == NULL) {
				map->array[index] = tmp->next;
			} else {
				parent->next = tmp->next;
			}	
			tmp->next = NULL;
			free(tmp);
			map->filled--;
			break;
		}
		parent = tmp;
	}
}
#endif
