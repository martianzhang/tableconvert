## Different Wiki Table Syntaxes

Explanation and examples of table-like syntax in **TWiki, MediaWiki, Confluence, TracWiki**:

### 1. **TWiki Syntax**  
TWiki uses its own table syntax with `|` as cell separators. Headers are bolded using `= Content =`, and row separation is done with `-` (optional but recommended below headers).  
```twiki
|=FIELD=|=TYPE=|=NULL=|=KEY=|=DEFAULT=|=EXTRA=|
|user_id|smallint(5)|NO|PRI|NULL|auto_increment|
|username|varchar(10)|NO||NULL||
|password|varchar(100)|NO||||
```


### 2. **MediaWiki Syntax**  
MediaWiki uses wiki table syntax starting with `{|` and ending with `|}`. Headers use `!`, data rows use `|`, and row separation uses `|-`.  
```mediawiki
{| class="wikitable" style="border: 1px solid #ddd;"
! FIELD !! TYPE !! NULL !! KEY !! DEFAULT !! EXTRA
|-
| user_id       | smallint(5) | NO   | PRI | NULL         | auto_increment
|-
| username      | varchar(10) | NO   |     | NULL         |
|-
| password      | varchar(100)| NO   |     |              |
|}
```


### 3. **Confluence Syntax**  
#### Method 1: Markdown Table (if Markdown support is enabled)  
When Confluence supports Markdown format, standard Markdown tables can be used directly:  
```markdown
| FIELD   | TYPE        | NULL | KEY | DEFAULT | EXTRA         |
|---------|-------------|------|-----|---------|---------------|
| user_id | smallint(5) | NO   | PRI | NULL    | auto_increment|
| username| varchar(10) | NO   |     | NULL    |               |
| password| varchar(100)| NO   |     |         |               |
```

#### Method 2: Native Confluence Table Syntax  
Cells are separated by `|`, and headers require no special markup (direct input, with drag-and-drop column resizing support):  
```plaintext
| FIELD   | TYPE        | NULL | KEY | DEFAULT | EXTRA         |
|---------|-------------|------|-----|---------|---------------|
| user_id | smallint(5) | NO   | PRI | NULL    | auto_increment|
| username| varchar(10) | NO   |     | NULL    |               |
| password| varchar(100)| NO   |     |         |               |
```


### 4. **TracWiki Syntax**  
TracWiki uses WikiText syntax, where headers are marked with `|= Content =|`, data rows use `| Content |`, and row separators use `---` (optional but recommended below headers).  
```tracwiki
|=FIELD=|=TYPE=|=NULL=|=KEY=|=DEFAULT=|=EXTRA=|
|user_id|smallint(5)|NO|PRI|NULL|auto_increment|
|username|varchar(10)|NO||NULL||
|password|varchar(100)|NO||||
```
