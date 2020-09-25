## PURPOSE


Many frameworks use ORM to interface with the database, GORM is typically the ORM used (also what I am using).
Yet there are always cases where you want to write the sql yourself. I am attempting to use the best of both worlds. Keep the ORM to constant lookups or simple ranges and then build individual SQL organized by the primary table.


## some musings 

dao - data access objects are also what you would call respository, but since respository sort of indicates a different database technology I'm using dao to indicate data access returned in generic *Rows type which is 

```
type Row struct {
    Data map[string]interface{}
    Len
}

type Rows []*Row

```

Then adapters by the caller is used

`var ent Entity  = entity.CovertRow(row)`

##
