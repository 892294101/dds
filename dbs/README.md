# dds


// protoc --go_out=. DataEventV1.proto
protoc *.proto --gofast_out=.

protoc *.proto --gofast_out=plugins=grpc:pcb



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


 pfile, err := spfile.LoadSpfile("D:\\workspace\\gowork\\src\\github.com/892294101\\dds\\build\\param\\httk_0001.desc",
		spfile.UTF8,
		log,
		spfile.GetMySQLName(),
		spfile.GetExtractName())
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := pfile.Production(); err != nil {
		log.Fatalf("%s", err)
	}


	ext := oramysql.NewMySQLSync()
	err = ext.InitSyncerConfig(log, pfile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = ext.StartSyncToStream(2, 1986)
	if err != nil {
		log.Fatalf("StartSyncToStream failed: %s", err)
	} 
	
	for i, row := range d.LCR.RowsEvent.Rows {

		for ind, value := range row {
			if value != nil {
				//d.ColValList[ind] = &protolib.ColumnValue{Existing: (*int32)(unsafe.Pointer(&value.Existing)), Length: (*int32)(unsafe.Pointer(&value.Length)), ColType: (*uint32)(unsafe.Pointer(&value.ColType)), PointField: (*int32)(unsafe.Pointer(&value.PointField)), ValBytes: value.ValBytes, ValInt32: &value.ValInt32, ValInt64: &value.ValInt64, ValInt8: (*int32)(unsafe.Pointer(&value.ValInt8)), ValInt16: (*int32)(unsafe.Pointer(&value.ValInt16)), ValString: &value.ValString, ValFloat32: &value.ValFloat32, ValFloat64: &value.ValFloat64, ValUint32: &value.ValUint32, ValUint64: &value.ValUint64, ValTime: timestamppb.New(value.ValTime), ValInt: (*int32)(unsafe.Pointer(&value.ValInt))}
				d.ColVal.Existing = (*int32)(unsafe.Pointer(&value.Existing))
				d.ColVal.Length = (*int32)(unsafe.Pointer(&value.Length))
				d.ColVal.ColType = (*uint32)(unsafe.Pointer(&value.ColType))
				d.ColVal.PointField = (*int32)(unsafe.Pointer(&value.PointField))
				d.ColVal.ValBytes = value.ValBytes
				d.ColVal.ValInt32 = &value.ValInt32
				d.ColVal.ValInt64 = &value.ValInt64
				d.ColVal.ValInt8 = (*int32)(unsafe.Pointer(&value.ValInt8))
				d.ColVal.ValInt16 = (*int32)(unsafe.Pointer(&value.ValInt16))
				d.ColVal.ValString = &value.ValString
				d.ColVal.ValFloat32 = &value.ValFloat32
				d.ColVal.ValFloat64 = &value.ValFloat64
				d.ColVal.ValUint32 = &value.ValUint32
				d.ColVal.ValUint64 = &value.ValUint64
				d.ColVal.ValTime = timestamppb.New(value.ValTime)
				d.ColVal.ValInt = (*int32)(unsafe.Pointer(&value.ValInt))
			}
			d.ColValList[ind] = d.ColVal
		}
		d.RowList[i] = &protolib.RowsList{Cv: d.ColValList}
	}	