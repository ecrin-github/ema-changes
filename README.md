# ema-changes
Identifies the Ids of trials recently updated by the EMA using an XML source file

Updates the source_data_studies table in the MDR with the date of the most recent file revision, or adds a record if the registry entry is new.
This data can thenm be used by the EU CTR downloading opprocess to obnly dopwenload the most recently revised / added files, on a weekly basis.
