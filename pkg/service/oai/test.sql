SELECT COUNT(*) AS num FROM coreview WHERE partition='mediathek' AND datestamp>='1970/01/01' AND (access='public' OR access='closed_data')


SELECT uuid, metadata, setspec, catalog, access, signature, sourcename, partition, status, seq, datestamp, identifier, url FROM coreview WHERE partition='mediathek' AND datestamp>='1970/01/01' AND (access='public' OR access='closed_data')


SELECT coreview.uuid FROM coreview LEFT JOIN pid ON coreview.uuid = pid.uuid WHERE coreview.partition='mediathek' AND pid.identifiertype='ARK' AND pid.uuid IS NULL LIMIT 100

SELECT coreview.uuid FROM coreview LEFT JOIN pid ON coreview.uuid = pid.uuid AND pid.identifiertype='ARK' WHERE coreview.partition='mediathek' AND pid.uuid IS NULL LIMIT 100