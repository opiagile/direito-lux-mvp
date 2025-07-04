-- Popular tabela processes com dados de teste realistas
-- Dados para todos os 8 tenants com processos variados

-- Silva & Associados (11111111-1111-1111-1111-111111111111)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('11111111-1111-1111-1111-111111111111', '5001234-12.2024.8.26.0100', 'TJSP', 'Ação de Cobrança - Contrato de Prestação de Serviços', 'active', true, NOW() - INTERVAL '30 days'),
('11111111-1111-1111-1111-111111111111', '5001235-45.2024.8.26.0224', 'TJSP', 'Ação Trabalhista - Rescisão Indireta', 'active', true, NOW() - INTERVAL '15 days'),
('11111111-1111-1111-1111-111111111111', '0011223-33.2024.5.02.0011', 'TRT2', 'Recurso Ordinário - Adicional de Insalubridade', 'active', false, NOW() - INTERVAL '45 days'),
('11111111-1111-1111-1111-111111111111', '5004567-89.2024.8.26.0053', 'TJSP', 'Execução de Título Extrajudicial', 'paused', false, NOW() - INTERVAL '60 days'),
('11111111-1111-1111-1111-111111111111', '5001111-22.2024.8.26.0100', 'TJSP', 'Ação de Despejo por Falta de Pagamento', 'active', true, NOW() - INTERVAL '7 days'),
('11111111-1111-1111-1111-111111111111', '1234567-89.2023.8.26.0100', 'TJSP', 'Ação de Indenização por Danos Morais', 'archived', false, NOW() - INTERVAL '200 days'),

-- Costa Santos Advogados (22222222-2222-2222-2222-222222222222)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('22222222-2222-2222-2222-222222222222', '5002001-11.2024.8.26.0224', 'TJSP', 'Ação de Divórcio Consensual', 'active', false, NOW() - INTERVAL '20 days'),
('22222222-2222-2222-2222-222222222222', '5002002-22.2024.8.26.0100', 'TJSP', 'Inventário e Partilha', 'active', true, NOW() - INTERVAL '35 days'),
('22222222-2222-2222-2222-222222222222', '5002003-33.2024.8.26.0053', 'TJSP', 'Ação de Alimentos', 'active', true, NOW() - INTERVAL '10 days'),
('22222222-2222-2222-2222-222222222222', '0022334-44.2024.5.02.0011', 'TRT2', 'Ação Trabalhista - Horas Extras', 'paused', false, NOW() - INTERVAL '50 days'),
('22222222-2222-2222-2222-222222222222', '5002004-55.2024.8.26.0224', 'TJSP', 'Usucapião Urbano', 'active', false, NOW() - INTERVAL '80 days'),

-- Barros Empresa (33333333-3333-3333-3333-333333333333)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('33333333-3333-3333-3333-333333333333', '5003001-11.2024.8.26.0053', 'TJSP', 'Ação de Cobrança - Prestação de Serviços Empresariais', 'active', true, NOW() - INTERVAL '25 days'),
('33333333-3333-3333-3333-333333333333', '5003002-22.2024.8.26.0100', 'TJSP', 'Dissolução de Sociedade Empresária', 'active', true, NOW() - INTERVAL '40 days'),
('33333333-3333-3333-3333-333333333333', '5003003-33.2024.8.26.0224', 'TJSP', 'Ação Anulatória de Débito Fiscal', 'active', false, NOW() - INTERVAL '15 days'),
('33333333-3333-3333-3333-333333333333', '1003004-44.2024.3.01.0001', 'TRF3', 'Mandado de Segurança - Tributos Federais', 'active', true, NOW() - INTERVAL '55 days'),

-- Lima Advogados (44444444-4444-4444-4444-444444444444)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('44444444-4444-4444-4444-444444444444', '5004001-11.2024.8.26.0100', 'TJSP', 'Ação de Revisão de Contrato Bancário', 'active', true, NOW() - INTERVAL '18 days'),
('44444444-4444-4444-4444-444444444444', '5004002-22.2024.8.26.0053', 'TJSP', 'Ação Consignatória em Pagamento', 'active', false, NOW() - INTERVAL '30 days'),
('44444444-4444-4444-4444-444444444444', '5004003-33.2024.8.26.0224', 'TJSP', 'Execução Fiscal - IPTU', 'paused', false, NOW() - INTERVAL '70 days'),
('44444444-4444-4444-4444-444444444444', '0044556-66.2024.5.02.0011', 'TRT2', 'Ação Trabalhista - Equiparação Salarial', 'active', true, NOW() - INTERVAL '12 days'),
('44444444-4444-4444-4444-444444444444', '5004004-77.2024.8.26.0100', 'TJSP', 'Ação de Reintegração de Posse', 'active', true, NOW() - INTERVAL '5 days'),

-- Pereira Advocacia (55555555-5555-5555-5555-555555555555)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('55555555-5555-5555-5555-555555555555', '5005001-11.2024.8.26.0053', 'TJSP', 'Ação de Responsabilidade Civil - Acidente de Trânsito', 'active', true, NOW() - INTERVAL '22 days'),
('55555555-5555-5555-5555-555555555555', '5005002-22.2024.8.26.0100', 'TJSP', 'Ação Declaratória de Inexistência de Débito', 'active', false, NOW() - INTERVAL '38 days'),
('55555555-5555-5555-5555-555555555555', '5005003-33.2024.8.26.0224', 'TJSP', 'Ação de Rescisão Contratual', 'active', true, NOW() - INTERVAL '8 days'),

-- Rodrigues Global (66666666-6666-6666-6666-666666666666)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('66666666-6666-6666-6666-666666666666', '5006001-11.2024.8.26.0100', 'TJSP', 'Ação de Propriedade Intelectual - Marca', 'active', true, NOW() - INTERVAL '28 days'),
('66666666-6666-6666-6666-666666666666', '5006002-22.2024.8.26.0053', 'TJSP', 'Ação de Concorrência Desleal', 'active', true, NOW() - INTERVAL '42 days'),
('66666666-6666-6666-6666-666666666666', '5006003-33.2024.8.26.0224', 'TJSP', 'Ação de Indenização - Responsabilidade Civil', 'active', false, NOW() - INTERVAL '14 days'),
('66666666-6666-6666-6666-666666666666', '1006004-44.2024.3.01.0001', 'TRF3', 'Ação Anulatória - Licitação Pública', 'active', true, NOW() - INTERVAL '60 days'),
('66666666-6666-6666-6666-666666666666', '5006005-55.2024.8.26.0100', 'TJSP', 'Arbitragem Empresarial', 'paused', false, NOW() - INTERVAL '90 days'),

-- Oliveira Partners (77777777-7777-7777-7777-777777777777)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('77777777-7777-7777-7777-777777777777', '5007001-11.2024.8.26.0224', 'TJSP', 'Ação de Recuperação Judicial', 'active', true, NOW() - INTERVAL '65 days'),
('77777777-7777-7777-7777-777777777777', '5007002-22.2024.8.26.0100', 'TJSP', 'Ação de Dissolução Irregular de Sociedade', 'active', true, NOW() - INTERVAL '35 days'),
('77777777-7777-7777-7777-777777777777', '5007003-33.2024.8.26.0053', 'TJSP', 'Ação de Cobrança - Duplicata', 'active', false, NOW() - INTERVAL '20 days'),
('77777777-7777-7777-7777-777777777777', '0077889-99.2024.5.02.0011', 'TRT2', 'Ação Trabalhista - Terceirização Ilícita', 'active', true, NOW() - INTERVAL '45 days'),

-- Machado Advogados (88888888-8888-8888-8888-888888888888)
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
('88888888-8888-8888-8888-888888888888', '5008001-11.2024.8.26.0100', 'TJSP', 'Ação Penal - Estelionato', 'active', true, NOW() - INTERVAL '33 days'),
('88888888-8888-8888-8888-888888888888', '5008002-22.2024.8.26.0053', 'TJSP', 'Habeas Corpus', 'active', true, NOW() - INTERVAL '12 days'),
('88888888-8888-8888-8888-888888888888', '5008003-33.2024.8.26.0224', 'TJSP', 'Ação de Reparação de Danos - Crime', 'active', false, NOW() - INTERVAL '25 days'),
('88888888-8888-8888-8888-888888888888', '1008004-44.2024.4.03.6100', 'TRF4', 'Ação Penal - Lavagem de Dinheiro', 'paused', true, NOW() - INTERVAL '120 days'),
('88888888-8888-8888-8888-888888888888', '5008005-55.2024.8.26.0100', 'TJSP', 'Revisão Criminal', 'active', false, NOW() - INTERVAL '6 days');

-- Adicionar alguns processos mais recentes (últimos 7 dias) para simular atividade
INSERT INTO processes (tenant_id, number, court, subject, status, monitoring, created_at) VALUES
-- Processos desta semana
('11111111-1111-1111-1111-111111111111', '5001999-88.2025.8.26.0100', 'TJSP', 'Ação de Cobrança - Janeiro 2025', 'active', true, NOW() - INTERVAL '3 days'),
('22222222-2222-2222-2222-222222222222', '5002999-88.2025.8.26.0224', 'TJSP', 'Divórcio Litigioso - Janeiro 2025', 'active', true, NOW() - INTERVAL '2 days'),
('33333333-3333-3333-3333-333333333333', '5003999-88.2025.8.26.0053', 'TJSP', 'Mandado de Segurança Empresarial', 'active', true, NOW() - INTERVAL '1 day'),
('44444444-4444-4444-4444-444444444444', '5004999-88.2025.8.26.0100', 'TJSP', 'Revisão de Contrato Bancário - 2025', 'active', true, NOW() - INTERVAL '4 days'),

-- Processos de hoje
('55555555-5555-5555-5555-555555555555', '5005999-88.2025.8.26.0053', 'TJSP', 'Ação de Indenização - Novo Cliente', 'active', true, NOW() - INTERVAL '3 hours'),
('66666666-6666-6666-6666-666666666666', '5006999-88.2025.8.26.0224', 'TJSP', 'Propriedade Intelectual - Software', 'active', true, NOW() - INTERVAL '1 hour');