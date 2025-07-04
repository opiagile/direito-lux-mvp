-- Popular tabela processes com dados de teste realistas
-- Adaptado para a estrutura existente da tabela

-- Limpar dados existentes de teste (se existirem)
DELETE FROM processes WHERE number LIKE '500%' OR number LIKE '001%' OR number LIKE '100%';

-- Silva & Associados (11111111-1111-1111-1111-111111111111)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('11111111-1111-1111-1111-111111111111', '5001234-12.2024.8.26.0100', 'Ação de Cobrança - Prestação de Serviços', 'Cobrança de honorários advocatícios e serviços jurídicos prestados ao cliente XYZ Ltda.', 'TJSP', 'active', true, NOW() - INTERVAL '30 days'),
('11111111-1111-1111-1111-111111111111', '5001235-45.2024.8.26.0224', 'Ação Trabalhista - Rescisão Indireta', 'Reclamação trabalhista pleiteando rescisão indireta por falta de pagamento de salários.', 'TJSP', 'active', true, NOW() - INTERVAL '15 days'),
('11111111-1111-1111-1111-111111111111', '0011223-33.2024.5.02.0011', 'Recurso Ordinário - Adicional de Insalubridade', 'Recurso contra sentença que negou adicional de insalubridade ao trabalhador.', 'TRT2', 'active', false, NOW() - INTERVAL '45 days'),
('11111111-1111-1111-1111-111111111111', '5004567-89.2024.8.26.0053', 'Execução de Título Extrajudicial', 'Execução de nota promissória no valor de R$ 50.000,00.', 'TJSP', 'paused', false, NOW() - INTERVAL '60 days'),
('11111111-1111-1111-1111-111111111111', '5001111-22.2024.8.26.0100', 'Ação de Despejo por Falta de Pagamento', 'Despejo por falta de pagamento de aluguéis em atraso há 6 meses.', 'TJSP', 'active', true, NOW() - INTERVAL '7 days'),
('11111111-1111-1111-1111-111111111111', '1234567-89.2023.8.26.0100', 'Ação de Indenização por Danos Morais', 'Indenização por danos morais decorrentes de negativação indevida.', 'TJSP', 'archived', false, NOW() - INTERVAL '200 days'),

-- Costa Santos Advogados (22222222-2222-2222-2222-222222222222)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('22222222-2222-2222-2222-222222222222', '5002001-11.2024.8.26.0224', 'Ação de Divórcio Consensual', 'Divórcio consensual com partilha de bens e guarda compartilhada.', 'TJSP', 'active', false, NOW() - INTERVAL '20 days'),
('22222222-2222-2222-2222-222222222222', '5002002-22.2024.8.26.0100', 'Inventário e Partilha', 'Inventário dos bens deixados pelo falecido Sr. João da Silva.', 'TJSP', 'active', true, NOW() - INTERVAL '35 days'),
('22222222-2222-2222-2222-222222222222', '5002003-33.2024.8.26.0053', 'Ação de Alimentos', 'Fixação de pensão alimentícia para menor de idade.', 'TJSP', 'active', true, NOW() - INTERVAL '10 days'),
('22222222-2222-2222-2222-222222222222', '0022334-44.2024.5.02.0011', 'Ação Trabalhista - Horas Extras', 'Reclamação trabalhista para pagamento de horas extras não quitadas.', 'TRT2', 'paused', false, NOW() - INTERVAL '50 days'),
('22222222-2222-2222-2222-222222222222', '5002004-55.2024.8.26.0224', 'Usucapião Urbano', 'Ação de usucapião de terreno urbano ocupado há mais de 15 anos.', 'TJSP', 'active', false, NOW() - INTERVAL '80 days'),

-- Barros Empresa (33333333-3333-3333-3333-333333333333)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('33333333-3333-3333-3333-333333333333', '5003001-11.2024.8.26.0053', 'Ação de Cobrança - Serviços Empresariais', 'Cobrança de serviços de consultoria empresarial prestados.', 'TJSP', 'active', true, NOW() - INTERVAL '25 days'),
('33333333-3333-3333-3333-333333333333', '5003002-22.2024.8.26.0100', 'Dissolução de Sociedade Empresária', 'Dissolução irregular de sociedade limitada com apuração de haveres.', 'TJSP', 'active', true, NOW() - INTERVAL '40 days'),
('33333333-3333-3333-3333-333333333333', '5003003-33.2024.8.26.0224', 'Ação Anulatória de Débito Fiscal', 'Anulação de auto de infração do município por cobrança indevida.', 'TJSP', 'active', false, NOW() - INTERVAL '15 days'),
('33333333-3333-3333-3333-333333333333', '1003004-44.2024.3.01.0001', 'Mandado de Segurança - Tributos Federais', 'MS contra cobrança de tributos federais considerados inconstitucionais.', 'TRF3', 'active', true, NOW() - INTERVAL '55 days'),

-- Lima Advogados (44444444-4444-4444-4444-444444444444)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('44444444-4444-4444-4444-444444444444', '5004001-11.2024.8.26.0100', 'Ação de Revisão de Contrato Bancário', 'Revisão de cláusulas abusivas em contrato de financiamento bancário.', 'TJSP', 'active', true, NOW() - INTERVAL '18 days'),
('44444444-4444-4444-4444-444444444444', '5004002-22.2024.8.26.0053', 'Ação Consignatória em Pagamento', 'Consignação em pagamento de valores contestados pelo credor.', 'TJSP', 'active', false, NOW() - INTERVAL '30 days'),
('44444444-4444-4444-4444-444444444444', '5004003-33.2024.8.26.0224', 'Execução Fiscal - IPTU', 'Execução fiscal de débitos de IPTU em atraso há 3 anos.', 'TJSP', 'paused', false, NOW() - INTERVAL '70 days'),
('44444444-4444-4444-4444-444444444444', '0044556-66.2024.5.02.0011', 'Ação Trabalhista - Equiparação Salarial', 'Ação pleiteando equiparação salarial entre empregados da mesma função.', 'TRT2', 'active', true, NOW() - INTERVAL '12 days'),
('44444444-4444-4444-4444-444444444444', '5004004-77.2024.8.26.0100', 'Ação de Reintegração de Posse', 'Reintegração de posse de imóvel rural invadido por terceiros.', 'TJSP', 'active', true, NOW() - INTERVAL '5 days'),

-- Pereira Advocacia (55555555-5555-5555-5555-555555555555)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('55555555-5555-5555-5555-555555555555', '5005001-11.2024.8.26.0053', 'Ação de Responsabilidade Civil - Acidente', 'Indenização por danos materiais e morais em acidente de trânsito.', 'TJSP', 'active', true, NOW() - INTERVAL '22 days'),
('55555555-5555-5555-5555-555555555555', '5005002-22.2024.8.26.0100', 'Ação Declaratória de Inexistência de Débito', 'Declaração de inexistência de débito cobrado indevidamente.', 'TJSP', 'active', false, NOW() - INTERVAL '38 days'),
('55555555-5555-5555-5555-555555555555', '5005003-33.2024.8.26.0224', 'Ação de Rescisão Contratual', 'Rescisão de contrato de prestação de serviços por inadimplemento.', 'TJSP', 'active', true, NOW() - INTERVAL '8 days'),

-- Rodrigues Global (66666666-6666-6666-6666-666666666666)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('66666666-6666-6666-6666-666666666666', '5006001-11.2024.8.26.0100', 'Ação de Propriedade Intelectual - Marca', 'Ação de nulidade de registro de marca por anterioridade.', 'TJSP', 'active', true, NOW() - INTERVAL '28 days'),
('66666666-6666-6666-6666-666666666666', '5006002-22.2024.8.26.0053', 'Ação de Concorrência Desleal', 'Ação por concorrência desleal e violação de segredo empresarial.', 'TJSP', 'active', true, NOW() - INTERVAL '42 days'),
('66666666-6666-6666-6666-666666666666', '5006003-33.2024.8.26.0224', 'Ação de Indenização - Responsabilidade Civil', 'Indenização por danos causados por produto defeituoso.', 'TJSP', 'active', false, NOW() - INTERVAL '14 days'),
('66666666-6666-6666-6666-666666666666', '1006004-44.2024.3.01.0001', 'Ação Anulatória - Licitação Pública', 'Anulação de licitação pública por vícios no edital.', 'TRF3', 'active', true, NOW() - INTERVAL '60 days'),
('66666666-6666-6666-6666-666666666666', '5006005-55.2024.8.26.0100', 'Arbitragem Empresarial', 'Procedimento arbitral para resolução de conflito societário.', 'TJSP', 'paused', false, NOW() - INTERVAL '90 days'),

-- Oliveira Partners (77777777-7777-7777-7777-777777777777)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('77777777-7777-7777-7777-777777777777', '5007001-11.2024.8.26.0224', 'Ação de Recuperação Judicial', 'Pedido de recuperação judicial de empresa em crise financeira.', 'TJSP', 'active', true, NOW() - INTERVAL '65 days'),
('77777777-7777-7777-7777-777777777777', '5007002-22.2024.8.26.0100', 'Ação de Dissolução Irregular de Sociedade', 'Dissolução irregular e apuração de haveres de sócio retirante.', 'TJSP', 'active', true, NOW() - INTERVAL '35 days'),
('77777777-7777-7777-7777-777777777777', '5007003-33.2024.8.26.0053', 'Ação de Cobrança - Duplicata', 'Cobrança de duplicatas mercantis vencidas e não pagas.', 'TJSP', 'active', false, NOW() - INTERVAL '20 days'),
('77777777-7777-7777-7777-777777777777', '0077889-99.2024.5.02.0011', 'Ação Trabalhista - Terceirização Ilícita', 'Reclamação trabalhista alegando terceirização ilícita de atividade-fim.', 'TRT2', 'active', true, NOW() - INTERVAL '45 days'),

-- Machado Advogados (88888888-8888-8888-8888-888888888888)
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
('88888888-8888-8888-8888-888888888888', '5008001-11.2024.8.26.0100', 'Ação Penal - Estelionato', 'Ação penal por crime de estelionato mediante fraude eletrônica.', 'TJSP', 'active', true, NOW() - INTERVAL '33 days'),
('88888888-8888-8888-8888-888888888888', '5008002-22.2024.8.26.0053', 'Habeas Corpus', 'HC preventivo para garantir direito de locomoção do paciente.', 'TJSP', 'active', true, NOW() - INTERVAL '12 days'),
('88888888-8888-8888-8888-888888888888', '5008003-33.2024.8.26.0224', 'Ação de Reparação de Danos - Crime', 'Ação de reparação civil por danos decorrentes de crime contra a honra.', 'TJSP', 'active', false, NOW() - INTERVAL '25 days'),
('88888888-8888-8888-8888-888888888888', '1008004-44.2024.4.03.6100', 'Ação Penal - Lavagem de Dinheiro', 'Ação penal por crime de lavagem de dinheiro em operação bancária.', 'TRF4', 'paused', true, NOW() - INTERVAL '120 days'),
('88888888-8888-8888-8888-888888888888', '5008005-55.2024.8.26.0100', 'Revisão Criminal', 'Revisão criminal para absolvição por erro judiciário.', 'TJSP', 'active', false, NOW() - INTERVAL '6 days');

-- Adicionar alguns processos mais recentes (últimos 7 dias) para simular atividade
INSERT INTO processes (tenant_id, number, title, description, court, status, monitoring_enabled, created_at) VALUES
-- Processos desta semana
('11111111-1111-1111-1111-111111111111', '5001999-88.2025.8.26.0100', 'Ação de Cobrança - Janeiro 2025', 'Nova ação de cobrança iniciada em janeiro de 2025.', 'TJSP', 'active', true, NOW() - INTERVAL '3 days'),
('22222222-2222-2222-2222-222222222222', '5002999-88.2025.8.26.0224', 'Divórcio Litigioso - Janeiro 2025', 'Processo de divórcio litigioso com disputa de bens.', 'TJSP', 'active', true, NOW() - INTERVAL '2 days'),
('33333333-3333-3333-3333-333333333333', '5003999-88.2025.8.26.0053', 'Mandado de Segurança Empresarial', 'MS contra ato administrativo que prejudica atividade empresarial.', 'TJSP', 'active', true, NOW() - INTERVAL '1 day'),
('44444444-4444-4444-4444-444444444444', '5004999-88.2025.8.26.0100', 'Revisão de Contrato Bancário - 2025', 'Nova revisão de contrato bancário com cláusulas abusivas.', 'TJSP', 'active', true, NOW() - INTERVAL '4 days'),

-- Processos de hoje
('55555555-5555-5555-5555-555555555555', '5005999-88.2025.8.26.0053', 'Ação de Indenização - Novo Cliente', 'Ação de indenização para cliente recém-contratado.', 'TJSP', 'active', true, NOW() - INTERVAL '3 hours'),
('66666666-6666-6666-6666-666666666666', '5006999-88.2025.8.26.0224', 'Propriedade Intelectual - Software', 'Ação de violação de direitos autorais de software.', 'TJSP', 'active', true, NOW() - INTERVAL '1 hour');