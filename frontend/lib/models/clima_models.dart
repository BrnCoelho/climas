class ClimaAtual {
  final int id;
  final String cidade;
  final String condicao;
  final int temperatura;
  final String vento;
  final int umidade;

  ClimaAtual({
    required this.id,
    required this.cidade,
    required this.condicao,
    required this.temperatura,
    required this.vento,
    required this.umidade,
  });

  factory ClimaAtual.fromJson(Map<String, dynamic> json) {
    return ClimaAtual(
      id: json['id'] ?? 0,
      cidade: json['cidade'] ?? 'N/A',
      condicao: json['condicao'] ?? 'N/A',
      temperatura: json['temperatura'] ?? 0,
      vento: json['vento'] ?? 'N/A',
      umidade: json['umidade'] ?? 0,
    );
  }
}

class PrevisaoDia {
  final String diaSemana;
  final int maxima;
  final int minima;
  final String condicao;

  PrevisaoDia({
    required this.diaSemana,
    required this.maxima,
    required this.minima,
    required this.condicao,
  });

  factory PrevisaoDia.fromJson(Map<String, dynamic> json) {
    return PrevisaoDia(
      diaSemana: json['dia_semana'] ?? 'N/A',
      maxima: json['maxima'] ?? 0,
      minima: json['minima'] ?? 0,
      condicao: json['condicao'] ?? 'N/A',
    );
  }
}