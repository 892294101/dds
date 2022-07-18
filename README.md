# dds

当前操作支持情况:
TABLE DDL(
    Rename Table                        --支持
    Alter Table                         --支持
    Drop Table                          --支持(另外还包含DROP VIEW/DROP GLOBAL TEMPORARY TABLE/DROP TEMPORARY TABLE)
                                              需要过滤掉(DROP GLOBAL TEMPORARY TABLE/DROP TEMPORARY TABLE)
    Create Table                        --支持
    Truncate Table                      --支持
)                                       
                                        
DATABASE DDL(                           
    CREATE DATABASE                     --支持
    ALTER DATABASE                      --支持
    DROP DATABASE                       --支持
)                                       
                                        
INDEX DDL(                              
    CREATE INDEX                        --支持
    DROP INDEX                          --支持
)

TEMPORARY DDL(
    CREATE TEMPORARY TABLE              --不支持
    DROP TEMPORARY TABLE                --不支持
    CREATE GLOBAL TEMPORARY TABLE       --不支持
    DROP GLOBAL TEMPORARY TABLE         --不支持
)

VIEW DDL(
    CREATE OR REPLACE VIEW              --支持
    DROP VIEW                           --支持,该事件归到了TABLE DDL
    ALTER VIEW                          --TODO
)                                       
                                        
INDEX DDL(                              
    CREATE SEQUENCE                     --支持
    DROP SEQUENCE                       --支持
    ALTER SEQUENCE                      --支持
)                                       
                                        
USER ROLE DDL(                          
    CREATE USER                         --支持
    ALTER USER                          --支持
    DROP USER                           --支持
    RENAME USER                         --支持
    GRANT PROXY                         --支持
    GRANT ROLE                          --支持
    GRANT STMT                          --支持
)

TRIGGER DDL(
    CREATE TRIGGER                      --TODO
    DROP TRIGGER                        --TODO
)                                       
                                        
PROCEDURE DDL(                          
    CREATE PROCEDURE                    --TODO
    DROP PROCEDURE                      --TODO
    ALTER PROCEDURE                     --TODO
)                                       
                                        
FUNCTION DDL(                           
    CREATE FUNCTION                     --TODO
    DROP FUNCTION                       --TODO
    ALTER FUNCTION                      --TODO
)

FUNCTION DDL(
    CREATE FUNCTION                     --TODO
    DROP FUNCTION                       --TODO
    ALTER FUNCTION                      --TODO
)


===============================================================================================
transaction	        `^SAVEPOINT`

skip all flush sqls	`^FLUSH`

table maintenance	`^OPTIMIZE\\s+TABLE`
                    `^ANALYZE\\s+TABLE`
                    `^REPAIR\\s+TABLE`
                    
temporary table	    `^DROP\\s+(\\/\\*\\!40005\\s+)?TEMPORARY\\s+(\\*\\/\\s+)?TABLE`

trigger	            `^CREATE\\s+(DEFINER\\s?=.+?)?TRIGGER`
                    `^DROP\\s+TRIGGER`
                    
procedure	        `^DROP\\s+PROCEDURE`
                    `^CREATE\\s+(DEFINER\\s?=.+?)?PROCEDURE`
                    `^ALTER\\s+PROCEDURE`
                    
view	            `^CREATE\\s*(OR REPLACE)?\\s+(ALGORITHM\\s?=.+?)?(DEFINER\\s?=.+?)?\\s+(SQL SECURITY DEFINER)?VIEW`
                    `^DROP\\s+VIEW`
                    `^ALTER\\s+(ALGORITHM\\s?=.+?)?(DEFINER\\s?=.+?)?(SQL SECURITY DEFINER)?VIEW`
                    
function	        `^CREATE\\s+(AGGREGATE)?\\s*?FUNCTION`
                    `^CREATE\\s+(DEFINER\\s?=.+?)?FUNCTION`
                    `^ALTER\\s+FUNCTION`
                    `^DROP\\s+FUNCTION`
                    
tableSpace	        `^CREATE\\s+TABLESPACE`
                    `^ALTER\\s+TABLESPACE`
                    `^DROP\\s+TABLESPACE`
                    
event	            `^CREATE\\s+(DEFINER\\s?=.+?)?EVENT`
                    `^ALTER\\s+(DEFINER\\s?=.+?)?EVENT`
                    `^DROP\\s+EVENT`
                    
account management	`^GRANT`
                    `^REVOKE`
                    `^CREATE\\s+USER`
                    `^ALTER\\s+USER`
                    `^RENAME\\s+USER`
                    `^DROP\\s+USER`
                    `^DROP\\s+USER`




root@localhost [admin]> desc geom6;
+-------+-------------------- 
| Field | Type                
+-------+-------------------- 
| n     | varchar(255)          
| g     | geometry              前4个字节存储srid,就是8位。SRID (4 bytes) + WKB。 spatial reference system identifier（SRID） Well-Known Binary（WKB）
| c     | geometrycollection  
| o     | point               
| e     | multipoint          
| a     | linestring          
| f     | multipolygon        
| d     | multilinestring     
| b     | polygon             
+-------+-------------------- 


insert into `admin`.`geom6` values(
'GIS 测试',
st_geomfromtext('point(50 70)'),
st_geomfromtext('geometrycollection(point(1 1),linestring(0 0,1 1,2 2,3 3,4 4))'),
st_geomfromtext('point(30 30)'),
st_mpointfromtext('multipoint ((1 1), (2 2), (3 3))') ,
st_geomfromtext('linestring(15 15, 20 20)'),
st_geomfromtext('multipolygon(((0 0,0 3,3 3,3 0,0 0),(1 1,1 2,2 2,2 1,1 1)))'),
st_geomfromtext('multilinestring((1 1,2 2,3 3),(4 4,5 5))'),
st_geomfromtext('polygon((0 0,0 3,3 0,0 0),(1 1,1 2,2 1,1 1))')
);

insert into `admin`.`geom6`(g) values( st_geomfromtext('point(50 70)') );
insert into `admin`.`geom6`(c) values( st_geomfromtext('geometrycollection(point(1 1),linestring(0 0,1 1,2 2,3 3,4 4))') );
insert into `admin`.`geom6`(o) values( st_geomfromtext('point(30 30)') );
insert into `admin`.`geom6`(e) values( st_mpointfromtext('multipoint ((1 1), (2 2), (3 3))') );
insert into `admin`.`geom6`(a) values( st_geomfromtext('linestring(15 15, 20 20)') );
insert into `admin`.`geom6`(f) values( st_geomfromtext('multipolygon(((0 0,0 3,3 3,3 0,0 0),(1 1,1 2,2 2,2 1,1 1)))') );
insert into `admin`.`geom6`(d) values( st_geomfromtext('multilinestring((1 1,2 2,3 3),(4 4,5 5))') );
insert into `admin`.`geom6`(b) values( st_geomfromtext('polygon((0 0,0 3,3 0,0 0),(1 1,1 2,2 1,1 1))') ); 


