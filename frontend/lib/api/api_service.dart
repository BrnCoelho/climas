import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/clima_models.dart';

// ATENÇÃO: COLOQUE O IP DA SUA API GO AQUI!
const String baseUrl = 'http://26.33.184.131:8080';

Future<List<ClimaAtual>> fetchClimasAtuais() async {
  final response = await http.get(Uri.parse('$baseUrl/clima/all'));
  if (response.statusCode == 200) {
    List jsonResponse = json.decode(response.body);
    return jsonResponse.map((data) => ClimaAtual.fromJson(data)).toList();
  } else {
    throw Exception('Falha ao carregar lista de climas');
  }
}

Future<List<PrevisaoDia>> fetchPrevisaoPorId(int idClimaAtual) async {
  final response =
      await http.get(Uri.parse('$baseUrl/previsao/por-clima/$idClimaAtual'));
  if (response.statusCode == 200) {
    List jsonResponse = json.decode(response.body);
    return jsonResponse.map((data) => PrevisaoDia.fromJson(data)).toList();
  } else {
    throw Exception('Falha ao carregar previsão do tempo');
  }
}