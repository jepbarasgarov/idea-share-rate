ALTER TABLE tbl_user_idea_rel ADD CONSTRAINT fk_idea FOREIGN KEY(idea_id) REFERENCES tbl_idea(id) ON DELETE CASCADE;
