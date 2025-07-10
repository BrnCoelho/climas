use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::time::{sleep, Duration};

#[derive(Debug, Deserialize)]
struct ApiResponse {
    results: Results,
}

#[derive(Debug, Deserialize)]
struct Results {
    city: String,
    date: String,
    time: String,
    temp: i32,
    humidity: u8,
    wind_speedy: String,
    description: String,
    forecast: Vec<Forecast>,
}

#[derive(Debug, Deserialize)]
struct Forecast {
    weekday: String,
    max: i32,
    min: i32,
    description: String,
}

#[derive(Debug, Serialize)]
struct ClimaAtualPost {
    cidade: String,
    data: String,
    hora: String,
    temperatura: i32,
    umidade: u8,
    vento: String,
    condicao: String,
}

#[derive(Debug, Serialize)]
struct PrevisaoPost {
    dia_semana: String,
    maxima: i32,
    minima: i32,
    condicao: String,
    id_clima_atual: i32,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let client = Client::new();
    let api_key = "b501ef0a";
    let modo_inicial = false; // true para GET->POST inicial, false para apenas GET->PUT(Atualizações)

    let cidades = vec![
        "Campinas",
        "São Paulo",
        "Rio de Janeiro",
        "Belo Horizonte",
        "Curitiba",
        "Itabira",
        "joão monlevade",
        "ipatinga",
        
    ];

    let mut ids_clima: HashMap<String, i32> = HashMap::new();
    let mut ids_previsao: HashMap<String, Vec<i32>> = HashMap::new();

    if modo_inicial {
        println!("--- MODO INICIAL: Criando registros com POST ---");

        for city_name in &cidades {
            println!("----- Criando registro para cidade: {}", city_name);
            let url = format!("https://api.hgbrasil.com/weather?key={}&city_name={}", api_key, city_name);
            let response = client.get(&url).send().await?;
            if !response.status().is_success() {
                println!("Erro ao buscar dados para {}: {:?}", city_name, response.text().await?);
                continue;
            }

            let api_response: ApiResponse = response.json().await?;

            let clima = ClimaAtualPost {
                cidade: api_response.results.city.clone(),
                data: format!("{}T00:00:00Z", converter_data(&api_response.results.date)),
                hora: format!("{}:00", api_response.results.time),
                temperatura: api_response.results.temp,
                umidade: api_response.results.humidity,
                vento: api_response.results.wind_speedy.clone(),
                condicao: api_response.results.description.clone(),
            };

            let clima_post_url = "http://00.00.000.000:8080/clima"; //O ip foi modificado 
            let clima_response = client.post(clima_post_url).json(&clima).send().await?;
            if !clima_response.status().is_success() {
                println!("Erro ao adicionar clima atual: {:?}", clima_response.text().await?);
                continue;
            }

            let clima_result: serde_json::Value = clima_response.json().await?;
            let id_clima_atual = clima_result["id"].as_i64().unwrap() as i32;
            ids_clima.insert(city_name.to_string(), id_clima_atual);

            let previsao_post_url = "http://26.33.184.131:8080/previsao";
            let mut previsoes_ids = Vec::new();

            for forecast in api_response.results.forecast.iter().take(6) {
                let previsao = PrevisaoPost {
                    dia_semana: forecast.weekday.clone(),
                    maxima: forecast.max,
                    minima: forecast.min,
                    condicao: forecast.description.clone(),
                    id_clima_atual,
                };

                let previsao_response = client.post(previsao_post_url).json(&previsao).send().await?;

                if previsao_response.status().is_success() {
                    let previsao_result: serde_json::Value = previsao_response.json().await?;
                    previsoes_ids.push(previsao_result["id"].as_i64().unwrap() as i32);
                } else {
                    println!("Erro ao adicionar previsão: {:?}", previsao_response.text().await?);
                }
            }
            ids_previsao.insert(city_name.to_string(), previsoes_ids);
        }
    } else {
        println!("--- MODO ATUALIZAÇÃO: Usando IDs fixos para PUT ---");
        ids_clima.insert("Campinas".into(), 16);
        ids_clima.insert("São Paulo".into(), 17); 
        ids_clima.insert("Rio de Janeiro".into(), 18);
        ids_clima.insert("Belo Horizonte".into(), 19);
        ids_clima.insert("Curitiba".into(), 20);
        ids_clima.insert("Itabira".into(), 21);
        ids_clima.insert("joão monlevade".into(), 22); 
        ids_clima.insert("ipatinga".into(), 23);  

        ids_previsao.insert("Campinas".into(), vec![91, 92, 93, 94, 95, 96]);
        ids_previsao.insert("São Paulo".into(), vec![97, 98, 99, 100, 101, 102]);
        ids_previsao.insert("Rio de Janeiro".into(), vec![103, 104, 105, 106, 107, 108]);
        ids_previsao.insert("Belo Horizonte".into(), vec![109, 110, 111, 112, 113, 114]);
        ids_previsao.insert("Curitiba".into(), vec![115, 116, 117, 118, 119, 120]);
        ids_previsao.insert("Itabira".into(), vec![121, 122, 123, 124, 125, 126]);
        ids_previsao.insert("joão monlevade".into(), vec![127, 128, 129, 130, 131, 132]);
        ids_previsao.insert("ipatinga".into(), vec![133, 134, 135, 136, 137, 138]);
    }

    println!("--- Iniciando atualizações periódicas ---");
    loop {
        sleep(Duration::from_secs(10)).await; // X segundos

        for cidade in &cidades {
            println!("----- Atualizando a cidade: {}", cidade);

            let url = format!("https://api.hgbrasil.com/weather?key={}&city_name={}", api_key, cidade);
            let response = client.get(&url).send().await?;
            if !response.status().is_success() {
                println!("Erro ao buscar dados para {}: {:?}", cidade, response.text().await?);
                continue;
            }

            let api_response: ApiResponse = response.json().await?;

            if let Some(&id_clima) = ids_clima.get(&cidade.to_string()) {
                let clima = ClimaAtualPost {
                    cidade: api_response.results.city.clone(),
                    data: format!("{}T00:00:00Z", converter_data(&api_response.results.date)),
                    hora: format!("{}:00", api_response.results.time),
                    temperatura: api_response.results.temp,
                    umidade: api_response.results.humidity,
                    vento: api_response.results.wind_speedy.clone(),
                    condicao: api_response.results.description.clone(),
                };

                let put_clima_url = format!("http://26.33.184.131:8080/clima/{}", id_clima);
                client.put(&put_clima_url).json(&clima).send().await?;

                if let Some(ids_prev) = ids_previsao.get(&cidade.to_string()) {
                    for (i, id_prev) in ids_prev.iter().enumerate() {
                        if let Some(forecast) = api_response.results.forecast.get(i) {
                            let previsao = PrevisaoPost {
                                dia_semana: forecast.weekday.clone(),
                                maxima: forecast.max,
                                minima: forecast.min,
                                condicao: forecast.description.clone(),
                                id_clima_atual: id_clima,
                            };

                            let put_prev_url = format!("http://26.33.184.131:8080/previsao/{}", id_prev);
                            client.put(&put_prev_url).json(&previsao).send().await?;
                        }
                    }
                }
            }

            println!("Aguardando para atualizar a próxima cidade...\n");
            
        }
        println!("Aguardando para a próxima atualização...\n");
    }
}

fn converter_data(data_br: &str) -> String {
    let partes: Vec<&str> = data_br.split('/').collect();
    if partes.len() != 3 {
        return String::from("1970-01-01");
    }
    format!("{}-{}-{}", partes[2], partes[1], partes[0])
}
