ALTER TABLE materials
    ADD CONSTRAINT unique_materials_name UNIQUE (name);

ALTER TABLE services
    ADD CONSTRAINT unique_services_name UNIQUE (name);