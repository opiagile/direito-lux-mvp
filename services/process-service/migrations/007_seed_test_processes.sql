-- Migration: 007_seed_test_processes.sql  
-- Cria processos de teste para validação de funcionalidades
-- Processos distribuídos entre os tenants de teste

-- ============================================================================
-- PROCESSOS DE TESTE POR TENANT
-- ============================================================================

-- STARTER TENANT 01 - Silva & Associados (5 processos)
INSERT INTO processes (id, tenant_id, number, type, subject, court, status, priority, monitoring, created_by, lawyer, estimated_value, tags, created_at, updated_at) VALUES
('proc-starter-01-001', 'tenant-starter-01', '5001234-20.2024.4.03.6109', 'Cível', 'Ação de Cobrança', 'TJSP - 1ª Vara Cível Central', 'active', 'high', true, 'user-starter-01-lawyer', 'Ricardo Silva', 25000.00, '["cobrança", "inadimplência"]', NOW(), NOW()),
('proc-starter-01-002', 'tenant-starter-01', '5009876-15.2024.4.03.6109', 'Trabalhista', 'Reclamação Trabalhista', 'TRT - 2ª Região - 5ª Vara', 'active', 'medium', false, 'user-starter-01-lawyer', 'Ricardo Silva', 15000.00, '["trabalhista", "horas-extras"]', NOW(), NOW()),
('proc-starter-01-003', 'tenant-starter-01', '5005555-30.2024.4.03.6109', 'Família', 'Divórcio Consensual', 'TJSP - Vara de Família', 'concluded', 'low', false, 'user-starter-01-lawyer', 'Ricardo Silva', 0, '["família", "divórcio"]', NOW(), NOW()),
('proc-starter-01-004', 'tenant-starter-01', '5002222-45.2024.4.03.6109', 'Cível', 'Indenização por Danos Morais', 'TJSP - 2ª Vara Cível', 'active', 'medium', true, 'user-starter-01-lawyer', 'Ricardo Silva', 50000.00, '["danos-morais", "indenização"]', NOW(), NOW()),
('proc-starter-01-005', 'tenant-starter-01', '5007777-88.2024.4.03.6109', 'Consumidor', 'Ação Consumidor vs Banco', 'TJSP - JEC Consumidor', 'suspended', 'low', false, 'user-starter-01-lawyer', 'Ricardo Silva', 8000.00, '["consumidor", "banco"]', NOW(), NOW()),

-- STARTER TENANT 02 - Oliveira Advocacia (3 processos - tenant trial)
('proc-starter-02-001', 'tenant-starter-02', '5011111-22.2024.4.03.6109', 'Cível', 'Ação de Despejo', 'TJSP - 3ª Vara Cível', 'active', 'high', true, 'user-starter-02-lawyer', 'Marco Oliveira', 120000.00, '["despejo", "locação"]', NOW(), NOW()),
('proc-starter-02-002', 'tenant-starter-02', '5033333-44.2024.4.03.6109', 'Tributário', 'Mandado de Segurança - IPTU', 'TJSP - Vara da Fazenda', 'active', 'medium', true, 'user-starter-02-lawyer', 'Marco Oliveira', 30000.00, '["tributário", "iptu"]', NOW(), NOW()),
('proc-starter-02-003', 'tenant-starter-02', '5055555-66.2024.4.03.6109', 'Criminal', 'Defesa - Estelionato', 'TJSP - 1ª Vara Criminal', 'active', 'high', true, 'user-starter-02-lawyer', 'Marco Oliveira', 0, '["criminal", "estelionato"]', NOW(), NOW()),

-- PROFESSIONAL TENANT 01 - Costa Santos (10 processos)
('proc-prof-01-001', 'tenant-prof-01', '5100001-11.2024.4.03.6109', 'Empresarial', 'Dissolução Societária', 'TJSP - Vara Empresarial', 'active', 'high', true, 'user-prof-01-lawyer', 'Felipe Costa', 500000.00, '["empresarial", "dissolução"]', NOW(), NOW()),
('proc-prof-01-002', 'tenant-prof-01', '5100002-22.2024.4.03.6109', 'Trabalhista', 'Ação Coletiva Trabalhista', 'TRT - 2ª Região - 1ª Vara', 'active', 'high', true, 'user-prof-01-lawyer', 'Felipe Costa', 1000000.00, '["trabalhista", "coletiva"]', NOW(), NOW()),
('proc-prof-01-003', 'tenant-prof-01', '5100003-33.2024.4.03.6109', 'Cível', 'Responsabilidade Civil Médica', 'TJSP - 4ª Vara Cível', 'active', 'medium', true, 'user-prof-01-lawyer', 'Felipe Costa', 200000.00, '["responsabilidade-civil", "médica"]', NOW(), NOW()),
('proc-prof-01-004', 'tenant-prof-01', '5100004-44.2024.4.03.6109', 'Imobiliário', 'Usucapião Urbano', 'TJSP - Vara de Registros', 'active', 'low', false, 'user-prof-01-lawyer', 'Felipe Costa', 800000.00, '["imobiliário", "usucapião"]', NOW(), NOW()),
('proc-prof-01-005', 'tenant-prof-01', '5100005-55.2024.4.03.6109', 'Consumidor', 'Ação vs Operadora Telecom', 'TJSP - JEC Consumidor', 'concluded', 'medium', false, 'user-prof-01-lawyer', 'Felipe Costa', 50000.00, '["consumidor", "telecom"]', NOW(), NOW()),
('proc-prof-01-006', 'tenant-prof-01', '5100006-66.2024.4.03.6109', 'Previdenciário', 'Aposentadoria por Invalidez', 'JEF - Seção Judiciária SP', 'active', 'medium', true, 'user-prof-01-lawyer', 'Felipe Costa', 0, '["previdenciário", "invalidez"]', NOW(), NOW()),
('proc-prof-01-007', 'tenant-prof-01', '5100007-77.2024.4.03.6109', 'Tributário', 'Anulatória de Débito Fiscal', 'TJSP - 1ª Vara Fazenda', 'active', 'high', true, 'user-prof-01-lawyer', 'Felipe Costa', 2000000.00, '["tributário", "anulatória"]', NOW(), NOW()),
('proc-prof-01-008', 'tenant-prof-01', '5100008-88.2024.4.03.6109', 'Administrativo', 'Mandado Segurança - Licitação', 'TJSP - Vara Fazenda', 'suspended', 'medium', false, 'user-prof-01-lawyer', 'Felipe Costa', 5000000.00, '["administrativo", "licitação"]', NOW(), NOW()),
('proc-prof-01-009', 'tenant-prof-01', '5100009-99.2024.4.03.6109', 'Ambiental', 'Ação Civil Pública - Poluição', 'TJSP - Vara Ambiental', 'active', 'high', true, 'user-prof-01-lawyer', 'Felipe Costa', 10000000.00, '["ambiental", "poluição"]', NOW(), NOW()),
('proc-prof-01-010', 'tenant-prof-01', '5100010-00.2024.4.03.6109', 'Penal', 'Habeas Corpus', 'STJ - Superior Tribunal', 'active', 'urgent', true, 'user-prof-01-lawyer', 'Felipe Costa', 0, '["penal", "habeas-corpus"]', NOW(), NOW()),

-- PROFESSIONAL TENANT 02 - Pereira Lima (8 processos)
('proc-prof-02-001', 'tenant-prof-02', '5200001-11.2024.4.03.6109', 'Contrato', 'Revisão Contratual Bancária', 'TJSP - 5ª Vara Cível', 'active', 'medium', true, 'user-prof-02-lawyer', 'Bruno Pereira', 150000.00, '["contrato", "bancário"]', NOW(), NOW()),
('proc-prof-02-002', 'tenant-prof-02', '5200002-22.2024.4.03.6109', 'Sucessões', 'Inventário e Partilha', 'TJSP - Vara Sucessões', 'active', 'low', false, 'user-prof-02-lawyer', 'Bruno Pereira', 3000000.00, '["sucessões", "inventário"]', NOW(), NOW()),
('proc-prof-02-003', 'tenant-prof-02', '5200003-33.2024.4.03.6109', 'Societário', 'Recuperação Judicial', 'TJSP - 1ª Vara Falências', 'active', 'urgent', true, 'user-prof-02-lawyer', 'Bruno Pereira', 50000000.00, '["societário", "recuperação"]', NOW(), NOW()),
('proc-prof-02-004', 'tenant-prof-02', '5200004-44.2024.4.03.6109', 'Trabalhista', 'Adicional Insalubridade', 'TRT - 2ª Região - 10ª Vara', 'concluded', 'medium', false, 'user-prof-02-lawyer', 'Bruno Pereira', 80000.00, '["trabalhista", "insalubridade"]', NOW(), NOW()),
('proc-prof-02-005', 'tenant-prof-02', '5200005-55.2024.4.03.6109', 'Intelectual', 'Violação Marca Registrada', 'TJSP - Vara Propriedade Intelectual', 'active', 'high', true, 'user-prof-02-lawyer', 'Bruno Pereira', 5000000.00, '["propriedade-intelectual", "marca"]', NOW(), NOW()),
('proc-prof-02-006', 'tenant-prof-02', '5200006-66.2024.4.03.6109', 'Internacional', 'Homologação Sentença Estrangeira', 'STJ - Corte Especial', 'active', 'medium', true, 'user-prof-02-lawyer', 'Bruno Pereira', 10000000.00, '["internacional", "homologação"]', NOW(), NOW()),
('proc-prof-02-007', 'tenant-prof-02', '5200007-77.2024.4.03.6109', 'Regulatório', 'Ação vs ANATEL', 'TRF - 3ª Região', 'active', 'high', true, 'user-prof-02-lawyer', 'Bruno Pereira', 20000000.00, '["regulatório", "anatel"]', NOW(), NOW()),
('proc-prof-02-008', 'tenant-prof-02', '5200008-88.2024.4.03.6109', 'Concorrencial', 'Investigação Cartel', 'TJSP - Vara Empresarial', 'suspended', 'urgent', false, 'user-prof-02-lawyer', 'Bruno Pereira', 100000000.00, '["concorrencial", "cartel"]', NOW(), NOW()),

-- BUSINESS TENANT 01 - Machado Advogados (15 processos)
('proc-biz-01-001', 'tenant-biz-01', '5300001-11.2024.4.03.6109', 'M&A', 'Due Diligence Aquisição', 'Arbitragem CCI', 'active', 'urgent', true, 'user-biz-01-lawyer', 'André Machado', 500000000.00, '["ma", "due-diligence"]', NOW(), NOW()),
('proc-biz-01-002', 'tenant-biz-01', '5300002-22.2024.4.03.6109', 'Compliance', 'Investigação Interna', 'Procedimento Interno', 'active', 'urgent', true, 'user-biz-01-lawyer', 'André Machado', 0, '["compliance", "investigação"]', NOW(), NOW()),
('proc-biz-01-003', 'tenant-biz-01', '5300003-33.2024.4.03.6109', 'Bancário', 'Financiamento Estruturado', 'TJSP - Vara Empresarial', 'active', 'high', true, 'user-biz-01-lawyer', 'André Machado', 1000000000.00, '["bancário", "financiamento"]', NOW(), NOW()),
('proc-biz-01-004', 'tenant-biz-01', '5300004-44.2024.4.03.6109', 'Energia', 'Contrato PPA Solar', 'ANEEL - Processo Administrativo', 'active', 'high', true, 'user-biz-01-lawyer', 'André Machado', 2000000000.00, '["energia", "ppa"]', NOW(), NOW()),
('proc-biz-01-005', 'tenant-biz-01', '5300005-55.2024.4.03.6109', 'Infraestrutura', 'Concessão Rodoviária', 'TCU - Tribunal de Contas', 'active', 'medium', true, 'user-biz-01-lawyer', 'André Machado', 15000000000.00, '["infraestrutura", "concessão"]', NOW(), NOW()),
('proc-biz-01-006', 'tenant-biz-01', '5300006-66.2024.4.03.6109', 'Mineração', 'Licenciamento Ambiental', 'IBAMA - Processo Administrativo', 'active', 'high', true, 'user-biz-01-lawyer', 'André Machado', 5000000000.00, '["mineração", "licenciamento"]', NOW(), NOW()),
('proc-biz-01-007', 'tenant-biz-01', '5300007-77.2024.4.03.6109', 'Tech', 'LGPD - Adequação Normativa', 'ANPD - Processo Administrativo', 'concluded', 'medium', false, 'user-biz-01-lawyer', 'André Machado', 0, '["tech", "lgpd"]', NOW(), NOW()),
('proc-biz-01-008', 'tenant-biz-01', '5300008-88.2024.4.03.6109', 'Telecomunicações', 'Licitação 5G', 'ANATEL - Licitação', 'active', 'urgent', true, 'user-biz-01-lawyer', 'André Machado', 50000000000.00, '["telecom", "5g"]', NOW(), NOW()),
('proc-biz-01-009', 'tenant-biz-01', '5300009-99.2024.4.03.6109', 'Aviação', 'Slot Aeroportuário', 'ANAC - Processo', 'active', 'medium', true, 'user-biz-01-lawyer', 'André Machado', 500000000.00, '["aviação", "slot"]', NOW(), NOW()),
('proc-biz-01-010', 'tenant-biz-01', '5300010-00.2024.4.03.6109', 'Farmacêutico', 'Registro ANVISA', 'ANVISA - Processo', 'active', 'high', true, 'user-biz-01-lawyer', 'André Machado', 1000000000.00, '["farmacêutico", "anvisa"]', NOW(), NOW()),
('proc-biz-01-011', 'tenant-biz-01', '5300011-11.2024.4.03.6109', 'Seguros', 'Resseguro Internacional', 'SUSEP - Processo', 'active', 'medium', false, 'user-biz-01-lawyer', 'André Machado', 10000000000.00, '["seguros", "resseguro"]', NOW(), NOW()),
('proc-biz-01-012', 'tenant-biz-01', '5300012-22.2024.4.03.6109', 'Petróleo', 'ANP - Concessão E&P', 'ANP - Licitação', 'active', 'urgent', true, 'user-biz-01-lawyer', 'André Machado', 100000000000.00, '["petróleo", "ep"]', NOW(), NOW()),
('proc-biz-01-013', 'tenant-biz-01', '5300013-33.2024.4.03.6109', 'Agronegócio', 'Certificação Orgânica', 'MAPA - Processo', 'concluded', 'low', false, 'user-biz-01-lawyer', 'André Machado', 500000000.00, '["agronegócio", "orgânico"]', NOW(), NOW()),
('proc-biz-01-014', 'tenant-biz-01', '5300014-44.2024.4.03.6109', 'Portuário', 'Arrendamento Porto', 'ANTAQ - Licitação', 'active', 'high', true, 'user-biz-01-lawyer', 'André Machado', 20000000000.00, '["portuário", "arrendamento"]', NOW(), NOW()),
('proc-biz-01-015', 'tenant-biz-01', '5300015-55.2024.4.03.6109', 'ESG', 'Certificação Sustentabilidade', 'Certificadora Privada', 'active', 'medium', true, 'user-biz-01-lawyer', 'André Machado', 0, '["esg", "sustentabilidade"]', NOW(), NOW()),

-- BUSINESS TENANT 02 - Tribunal Consultoria (12 processos)
('proc-biz-02-001', 'tenant-biz-02', '5400001-11.2024.4.03.6109', 'Público', 'Consultoria Legislativa', 'Congresso Nacional', 'active', 'high', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["público", "legislativo"]', NOW(), NOW()),
('proc-biz-02-002', 'tenant-biz-02', '5400002-22.2024.4.03.6109', 'Constitucional', 'ADPF no STF', 'STF - Supremo Tribunal', 'active', 'urgent', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["constitucional", "adpf"]', NOW(), NOW()),
('proc-biz-02-003', 'tenant-biz-02', '5400003-33.2024.4.03.6109', 'Eleitoral', 'Ação Impugnação Mandato', 'TSE - Tribunal Superior', 'active', 'urgent', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["eleitoral", "mandato"]', NOW(), NOW()),
('proc-biz-02-004', 'tenant-biz-02', '5400004-44.2024.4.03.6109', 'Militar', 'Conselho Justiça Militar', 'STM - Superior Tribunal', 'active', 'medium', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["militar", "conselho"]', NOW(), NOW()),
('proc-biz-02-005', 'tenant-biz-02', '5400005-55.2024.4.03.6109', 'Trabalho', 'Dissídio Coletivo', 'TST - Tribunal Superior', 'active', 'high', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 50000000000.00, '["trabalho", "dissídio"]', NOW(), NOW()),
('proc-biz-02-006', 'tenant-biz-02', '5400006-66.2024.4.03.6109', 'Previdência', 'Uniformização Jurisprudência', 'TNU - Turma Nacional', 'active', 'medium', false, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["previdência", "uniformização"]', NOW(), NOW()),
('proc-biz-02-007', 'tenant-biz-02', '5400007-77.2024.4.03.6109', 'Administrativo', 'Controle Concentrado', 'STF - Pleno', 'active', 'urgent', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["administrativo", "concentrado"]', NOW(), NOW()),
('proc-biz-02-008', 'tenant-biz-02', '5400008-88.2024.4.03.6109', 'Internacional', 'Corte Interamericana DH', 'CIDH - Corte Internacional', 'active', 'high', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["internacional", "direitos-humanos"]', NOW(), NOW()),
('proc-biz-02-009', 'tenant-biz-02', '5400009-99.2024.4.03.6109', 'Ambiental', 'Mudanças Climáticas', 'STF - Ação Climática', 'active', 'urgent', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["ambiental", "clima"]', NOW(), NOW()),
('proc-biz-02-010', 'tenant-biz-02', '5400010-00.2024.4.03.6109', 'Digital', 'Marco Civil Internet', 'STF - ADI', 'concluded', 'medium', false, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["digital", "marco-civil"]', NOW(), NOW()),
('proc-biz-02-011', 'tenant-biz-02', '5400011-11.2024.4.03.6109', 'Bioética', 'Pesquisa Células-Tronco', 'STF - ADPF', 'active', 'medium', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 0, '["bioética", "células-tronco"]', NOW(), NOW()),
('proc-biz-02-012', 'tenant-biz-02', '5400012-22.2024.4.03.6109', 'Tributário', 'Modulação Efeitos STF', 'STF - RE com Repercussão', 'active', 'high', true, 'user-biz-02-lawyer', 'Gustavo Tribunal', 500000000000.00, '["tributário", "modulação"]', NOW(), NOW()),

-- ENTERPRISE TENANT 01 - Barros Enterprise (20 processos)  
('proc-ent-01-001', 'tenant-ent-01', '5500001-11.2024.4.03.6109', 'Corporate', 'Fusão Multinacional', 'Multiple Jurisdictions', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 500000000000.00, '["corporate", "fusão"]', NOW(), NOW()),
('proc-ent-01-002', 'tenant-ent-01', '5500002-22.2024.4.03.6109', 'Cross-border', 'Arbitragem ICC Paris', 'ICC - Arbitration Court', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 100000000000.00, '["arbitragem", "icc"]', NOW(), NOW()),
('proc-ent-01-003', 'tenant-ent-01', '5500003-33.2024.4.03.6109', 'Antitrust', 'Merger Control Global', 'Multiple Antitrust Authorities', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["antitrust", "merger"]', NOW(), NOW()),
('proc-ent-01-004', 'tenant-ent-01', '5500004-44.2024.4.03.6109', 'Financial', 'Derivative Master Agreement', 'ISDA - London', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 1000000000000.00, '["financial", "derivatives"]', NOW(), NOW()),
('proc-ent-01-005', 'tenant-ent-01', '5500005-55.2024.4.03.6109', 'Sovereign', 'Sovereign Bond Issuance', 'Multiple Markets', 'concluded', 'medium', false, 'user-ent-01-lawyer', 'Rodrigo Barros', 50000000000000.00, '["sovereign", "bonds"]', NOW(), NOW()),
('proc-ent-01-006', 'tenant-ent-01', '5500006-66.2024.4.03.6109', 'Infrastructure', 'PPP Megaproject', 'Brazilian Government', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 200000000000.00, '["infrastructure", "ppp"]', NOW(), NOW()),
('proc-ent-01-007', 'tenant-ent-01', '5500007-77.2024.4.03.6109', 'Technology', 'AI Governance Framework', 'Global Regulators', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["technology", "ai"]', NOW(), NOW()),
('proc-ent-01-008', 'tenant-ent-01', '5500008-88.2024.4.03.6109', 'ESG', 'Climate Litigation', 'Multiple Courts', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["esg", "climate"]', NOW(), NOW()),
('proc-ent-01-009', 'tenant-ent-01', '5500009-99.2024.4.03.6109', 'Space Law', 'Satellite Constellation', 'ITU - International', 'active', 'medium', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 100000000000.00, '["space", "satellite"]', NOW(), NOW()),
('proc-ent-01-010', 'tenant-ent-01', '5500010-00.2024.4.03.6109', 'Quantum', 'Quantum Computing IP', 'USPTO + EPO', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 50000000000.00, '["quantum", "ip"]', NOW(), NOW()),
('proc-ent-01-011', 'tenant-ent-01', '5500011-11.2024.4.03.6109', 'Biotech', 'Gene Therapy Approval', 'FDA + EMA + ANVISA', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 200000000000.00, '["biotech", "gene-therapy"]', NOW(), NOW()),
('proc-ent-01-012', 'tenant-ent-01', '5500012-22.2024.4.03.6109', 'Renewable', 'Offshore Wind Farm', 'Multiple Jurisdictions', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 500000000000.00, '["renewable", "offshore"]', NOW(), NOW()),
('proc-ent-01-013', 'tenant-ent-01', '5500013-33.2024.4.03.6109', 'Crypto', 'Digital Asset Regulation', 'Global Regulators', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["crypto", "regulation"]', NOW(), NOW()),
('proc-ent-01-014', 'tenant-ent-01', '5500014-44.2024.4.03.6109', 'Gaming', 'Metaverse Legal Framework', 'Multiple Jurisdictions', 'active', 'medium', false, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["gaming", "metaverse"]', NOW(), NOW()),
('proc-ent-01-015', 'tenant-ent-01', '5500015-55.2024.4.03.6109', 'Aviation', 'eVTOL Certification', 'FAA + EASA + ANAC', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 100000000000.00, '["aviation", "evtol"]', NOW(), NOW()),
('proc-ent-01-016', 'tenant-ent-01', '5500016-66.2024.4.03.6109', 'Maritime', 'Autonomous Ships Legal', 'IMO - International', 'active', 'medium', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 50000000000.00, '["maritime", "autonomous"]', NOW(), NOW()),
('proc-ent-01-017', 'tenant-ent-01', '5500017-77.2024.4.03.6109', 'Nuclear', 'Small Modular Reactor', 'IAEA + NRC + CNEN', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 1000000000000.00, '["nuclear", "smr"]', NOW(), NOW()),
('proc-ent-01-018', 'tenant-ent-01', '5500018-88.2024.4.03.6109', 'Cyber', 'Nation-State Attribution', 'International Courts', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["cyber", "attribution"]', NOW(), NOW()),
('proc-ent-01-019', 'tenant-ent-01', '5500019-99.2024.4.03.6109', 'Data', 'Global Data Governance', 'Multiple DPAs', 'active', 'high', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["data", "governance"]', NOW(), NOW()),
('proc-ent-01-020', 'tenant-ent-01', '5500020-00.2024.4.03.6109', 'Sovereign', 'Digital Currency CBDC', 'Central Banks Consortium', 'active', 'urgent', true, 'user-ent-01-lawyer', 'Rodrigo Barros', 0, '["sovereign", "cbdc"]', NOW(), NOW()),

-- ENTERPRISE TENANT 02 - Mega Advocacia (18 processos)
('proc-ent-02-001', 'tenant-ent-02', '5600001-11.2024.4.03.6109', 'Global M&A', 'Mega Merger Deal', 'Global Antitrust Review', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 2000000000000.00, '["mega-merger", "global"]', NOW(), NOW()),
('proc-ent-02-002', 'tenant-ent-02', '5600002-22.2024.4.03.6109', 'Pharmaceutical', 'Global Drug Launch', 'Worldwide Regulatory', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 500000000000.00, '["pharma", "global-launch"]', NOW(), NOW()),
('proc-ent-02-003', 'tenant-ent-02', '5600003-33.2024.4.03.6109', 'Technology', 'Platform Regulation', 'Big Tech Compliance', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 0, '["big-tech", "platform"]', NOW(), NOW()),
('proc-ent-02-004', 'tenant-ent-02', '5600004-44.2024.4.03.6109', 'Energy Transition', 'Green Hydrogen Economy', 'Global Energy Transition', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 1000000000000.00, '["hydrogen", "transition"]', NOW(), NOW()),
('proc-ent-02-005', 'tenant-ent-02', '5600005-55.2024.4.03.6109', 'Space Economy', 'Commercial Space Station', 'International Space Law', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 200000000000.00, '["space-economy", "station"]', NOW(), NOW()),
('proc-ent-02-006', 'tenant-ent-02', '5600006-66.2024.4.03.6109', 'AI Ethics', 'AI Safety Standards', 'Global AI Governance', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 0, '["ai-ethics", "safety"]', NOW(), NOW()),
('proc-ent-02-007', 'tenant-ent-02', '5600007-77.2024.4.03.6109', 'Quantum Internet', 'Quantum Communication', 'Quantum Infrastructure', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 1000000000000.00, '["quantum-internet", "communication"]', NOW(), NOW()),
('proc-ent-02-008', 'tenant-ent-02', '5600008-88.2024.4.03.6109', 'Synthetic Biology', 'Engineered Organisms', 'Biosafety Regulation', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 500000000000.00, '["synthetic-bio", "engineered"]', NOW(), NOW()),
('proc-ent-02-009', 'tenant-ent-02', '5600009-99.2024.4.03.6109', 'Climate Engineering', 'Geoengineering Governance', 'Global Climate Response', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 0, '["geoengineering", "climate"]', NOW(), NOW()),
('proc-ent-02-010', 'tenant-ent-02', '5600010-00.2024.4.03.6109', 'Neural Interface', 'Brain-Computer Interface', 'Neurotechnology Ethics', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 1000000000000.00, '["neural", "bci"]', NOW(), NOW()),
('proc-ent-02-011', 'tenant-ent-02', '5600011-11.2024.4.03.6109', 'Digital Twin Earth', 'Planet Simulation', 'Global Data Integration', 'active', 'medium', true, 'user-ent-02-lawyer', 'Thiago Mega', 2000000000000.00, '["digital-twin", "earth"]', NOW(), NOW()),
('proc-ent-02-012', 'tenant-ent-02', '5600012-22.2024.4.03.6109', 'Longevity Economy', 'Life Extension Tech', 'Bioethics Framework', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 500000000000.00, '["longevity", "bioethics"]', NOW(), NOW()),
('proc-ent-02-013', 'tenant-ent-02', '5600013-33.2024.4.03.6109', 'Ocean Economy', 'Deep Sea Mining', 'International Seabed Authority', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 1000000000000.00, '["ocean", "deep-sea"]', NOW(), NOW()),
('proc-ent-02-014', 'tenant-ent-02', '5600014-44.2024.4.03.6109', 'Asteroid Mining', 'Space Resource Rights', 'Outer Space Treaty', 'active', 'medium', false, 'user-ent-02-lawyer', 'Thiago Mega', 10000000000000.00, '["asteroid", "mining"]', NOW(), NOW()),
('proc-ent-02-015', 'tenant-ent-02', '5600015-55.2024.4.03.6109', 'Consciousness AI', 'AGI Rights Framework', 'AI Personhood Law', 'active', 'urgent', true, 'user-ent-02-lawyer', 'Thiago Mega', 0, '["agi", "consciousness"]', NOW(), NOW()),
('proc-ent-02-016', 'tenant-ent-02', '5600016-66.2024.4.03.6109', 'Molecular Assembly', 'Nanotechnology Safety', 'Molecular Manufacturing', 'active', 'high', true, 'user-ent-02-lawyer', 'Thiago Mega', 500000000000.00, '["nano", "molecular"]', NOW(), NOW()),
('proc-ent-02-017', 'tenant-ent-02', '5600017-77.2024.4.03.6109', 'Time Crystals', 'Quantum Materials', 'Novel Physics Applications', 'concluded', 'medium', false, 'user-ent-02-lawyer', 'Thiago Mega', 100000000000.00, '["time-crystals", "quantum"]', NOW(), NOW()),
('proc-ent-02-018', 'tenant-ent-02', '5600018-88.2024.4.03.6109', 'Multiverse Theory', 'Parallel Universe Law', 'Theoretical Legal Framework', 'active', 'low', false, 'user-ent-02-lawyer', 'Thiago Mega', 0, '["multiverse", "theoretical"]', NOW(), NOW());

-- ============================================================================
-- PARTES PROCESSUAIS (algumas para demonstração)
-- ============================================================================

-- Partes para alguns processos de exemplo
INSERT INTO parties (id, process_id, name, document, type, role, created_at, updated_at) VALUES
-- Processo Starter 01
('party-001', 'proc-starter-01-001', 'João da Silva Santos', '123.456.789-01', 'person', 'plaintiff', NOW(), NOW()),
('party-002', 'proc-starter-01-001', 'Empresa ABC Ltda', '12.345.678/0001-90', 'company', 'defendant', NOW(), NOW()),

-- Processo Professional 01  
('party-003', 'proc-prof-01-001', 'Mega Corporation S.A.', '11.111.111/0001-11', 'company', 'plaintiff', NOW(), NOW()),
('party-004', 'proc-prof-01-001', 'Super Empresa Ltda', '22.222.222/0001-22', 'company', 'defendant', NOW(), NOW()),

-- Processo Business 01
('party-005', 'proc-biz-01-001', 'Global Tech Corp', '33.333.333/0001-33', 'company', 'plaintiff', NOW(), NOW()),
('party-006', 'proc-biz-01-001', 'International Holdings Ltd', '44.444.444/0001-44', 'company', 'defendant', NOW(), NOW()),

-- Processo Enterprise 01
('party-007', 'proc-ent-01-001', 'Multinational Consortium', '55.555.555/0001-55', 'company', 'plaintiff', NOW(), NOW()),
('party-008', 'proc-ent-01-001', 'Global Investment Fund', '66.666.666/0001-66', 'company', 'defendant', NOW(), NOW());

-- ============================================================================
-- MOVIMENTAÇÕES PROCESSUAIS (algumas para demonstração)
-- ============================================================================

-- Movimentações para processos ativos
INSERT INTO movements (id, process_id, date, description, type, is_important, created_at, updated_at) VALUES
-- Processo Starter - movimentações simples
('mov-001', 'proc-starter-01-001', NOW() - INTERVAL '2 days', 'Petição inicial distribuída', 'distribuição', true, NOW(), NOW()),
('mov-002', 'proc-starter-01-001', NOW() - INTERVAL '1 day', 'Despacho inicial - cite-se', 'despacho', true, NOW(), NOW()),
('mov-003', 'proc-starter-01-001', NOW(), 'Mandado de citação expedido', 'expedição', false, NOW(), NOW()),

-- Processo Professional - movimentações médias
('mov-004', 'proc-prof-01-001', NOW() - INTERVAL '10 days', 'Petição inicial protocolada', 'protocolo', true, NOW(), NOW()),
('mov-005', 'proc-prof-01-001', NOW() - INTERVAL '8 days', 'Distribuição por sorteio', 'distribuição', true, NOW(), NOW()),
('mov-006', 'proc-prof-01-001', NOW() - INTERVAL '5 days', 'Despacho determinando emenda à inicial', 'despacho', true, NOW(), NOW()),
('mov-007', 'proc-prof-01-001', NOW() - INTERVAL '3 days', 'Petição de emenda protocolada', 'protocolo', false, NOW(), NOW()),
('mov-008', 'proc-prof-01-001', NOW() - INTERVAL '1 day', 'Inicial recebida - cite-se o réu', 'despacho', true, NOW(), NOW()),

-- Processo Business - movimentações complexas
('mov-009', 'proc-biz-01-001', NOW() - INTERVAL '30 days', 'Petição inicial de due diligence', 'protocolo', true, NOW(), NOW()),
('mov-010', 'proc-biz-01-001', NOW() - INTERVAL '25 days', 'Decisão liminar deferida', 'decisão', true, NOW(), NOW()),
('mov-011', 'proc-biz-01-001', NOW() - INTERVAL '20 days', 'Agravo de instrumento interposto', 'recurso', true, NOW(), NOW()),
('mov-012', 'proc-biz-01-001', NOW() - INTERVAL '15 days', 'Contraminuta do agravo', 'protocolo', false, NOW(), NOW()),
('mov-013', 'proc-biz-01-001', NOW() - INTERVAL '10 days', 'Agravo distribuído no TJ', 'distribuição', false, NOW(), NOW()),
('mov-014', 'proc-biz-01-001', NOW() - INTERVAL '5 days', 'Decisão monocrática - negou provimento', 'decisão', true, NOW(), NOW()),
('mov-015', 'proc-biz-01-001', NOW() - INTERVAL '2 days', 'Trânsito em julgado da liminar', 'trânsito', true, NOW(), NOW()),

-- Processo Enterprise - movimentações globais
('mov-016', 'proc-ent-01-001', NOW() - INTERVAL '60 days', 'International arbitration filed', 'filing', true, NOW(), NOW()),
('mov-017', 'proc-ent-01-001', NOW() - INTERVAL '50 days', 'Tribunal constituted', 'constitution', true, NOW(), NOW()),
('mov-018', 'proc-ent-01-001', NOW() - INTERVAL '40 days', 'Preliminary objections filed', 'objection', true, NOW(), NOW()),
('mov-019', 'proc-ent-01-001', NOW() - INTERVAL '30 days', 'Document production ordered', 'order', true, NOW(), NOW()),
('mov-020', 'proc-ent-01-001', NOW() - INTERVAL '20 days', 'Expert witnesses appointed', 'appointment', false, NOW(), NOW()),
('mov-021', 'proc-ent-01-001', NOW() - INTERVAL '10 days', 'Hearing scheduled', 'scheduling', true, NOW(), NOW()),
('mov-022', 'proc-ent-01-001', NOW() - INTERVAL '5 days', 'Pre-hearing briefs submitted', 'submission', false, NOW(), NOW()),
('mov-023', 'proc-ent-01-001', NOW() - INTERVAL '1 day', 'Evidentiary hearing held', 'hearing', true, NOW(), NOW());

-- ============================================================================
-- COMENTÁRIOS PARA ORIENTAÇÃO
-- ============================================================================

-- DADOS CRIADOS:
-- - 8 tenants (2 por plano)
-- - 32 usuários (4 roles por tenant) 
-- - ~90 processos distribuídos realisticamente por plano
-- - Partes e movimentações de exemplo
-- 
-- DISTRIBUIÇÃO DE PROCESSOS:
-- Starter: 5 + 3 = 8 processos (dentro do limite de 50)
-- Professional: 10 + 8 = 18 processos (dentro do limite de 200)  
-- Business: 15 + 12 = 27 processos (dentro do limite de 500)
-- Enterprise: 20 + 18 = 38 processos (ilimitado)
--
-- VARIAÇÃO REALÍSTICA:
-- - Processos simples para Starter (cobrança, família, etc.)
-- - Processos médios para Professional (empresarial, trabalhista)
-- - Processos complexos para Business (M&A, regulatório, infraestrutura)
-- - Processos globais para Enterprise (arbitragem internacional, tech avançada)
--
-- TODOS OS CAMPOS TESTÁVEIS:
-- - Status: active, concluded, suspended
-- - Priority: low, medium, high, urgent  
-- - Monitoring: true/false
-- - Tags: múltiplas áreas jurídicas
-- - Valores: desde R$ 0 até trilhões
-- - Courts: desde JEC até tribunais internacionais