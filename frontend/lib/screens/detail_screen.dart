import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:weather_icons/weather_icons.dart';
import '../api/api_service.dart';
import '../models/clima_models.dart';
import '../utils/weather_helper.dart';

class DetailScreen extends StatefulWidget {
  final ClimaAtual clima;

  const DetailScreen({super.key, required this.clima});

  @override
  State<DetailScreen> createState() => _DetailScreenState();
}

class _DetailScreenState extends State<DetailScreen> {
  late Future<List<PrevisaoDia>> futurePrevisao;

  @override
  void initState() {
    super.initState();
    futurePrevisao = fetchPrevisaoPorId(widget.clima.id);
  }

  @override
  Widget build(BuildContext context) {
    final String nomeCidadeCorrigido =
        widget.clima.cidade.toLowerCase() == 'sao paulo'
            ? 'São Paulo'
            : widget.clima.cidade;

    return Scaffold(
      appBar: AppBar(
        title: Text(nomeCidadeCorrigido),
      ),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20.0),
        child: Column(
          children: [
            const SizedBox(height: 20),
            // --- Seção do Clima Atual ---
            Text(
              nomeCidadeCorrigido,
              style: GoogleFonts.lato(
                  fontSize: 34,
                  color: Colors.white,
                  fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 10),
            Text(
              '${widget.clima.temperatura}°',
              style: GoogleFonts.lato(
                  fontSize: 80,
                  color: Colors.white,
                  fontWeight: FontWeight.w300),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                BoxedIcon(getWeatherIcon(widget.clima.condicao),
                    size: 30, color: Colors.white70),
                const SizedBox(width: 10),
                Text(
                  widget.clima.condicao,
                  style: GoogleFonts.lato(
                      fontSize: 20,
                      color: Colors.white70,
                      fontWeight: FontWeight.w500),
                ),
              ],
            ),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text('Vento: ${widget.clima.vento}',
                    style:
                        GoogleFonts.lato(fontSize: 16, color: Colors.white70)),
                const SizedBox(width: 20),
                Text('Umidade: ${widget.clima.umidade}%',
                    style:
                        GoogleFonts.lato(fontSize: 16, color: Colors.white70)),
              ],
            ),
            const SizedBox(height: 40),

            // --- Seção da Previsão ---
            const Text("Previsão para a Semana",
                style: TextStyle(
                    fontSize: 18,
                    color: Colors.white,
                    fontWeight: FontWeight.bold)),
            const SizedBox(height: 15),

            // O Expanded foi substituído por um SizedBox para controlar a altura.
            SizedBox(
              height:
                  125, // Altura reduzida para os cards. Você pode ajustar este valor.
              child: FutureBuilder<List<PrevisaoDia>>(
                future: futurePrevisao,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const Center(
                        child: CircularProgressIndicator(color: Colors.white));
                  } else if (snapshot.hasError) {
                    return Center(
                        child: Text('Não foi possível carregar a previsão.',
                            style: TextStyle(color: Colors.red[300])));
                  } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
                    final previsoes = snapshot.data!.take(6).toList();
                    // A ListView agora está dentro de um Center para alinhar o conjunto no meio.
                    return Center(
                      child: ListView.separated(
                        scrollDirection: Axis.horizontal,
                        shrinkWrap:
                            true, // Permite que a lista tenha o tamanho do seu conteúdo.
                        itemCount: previsoes.length,
                        separatorBuilder: (context, index) =>
                            const SizedBox(width: 12),
                        itemBuilder: (context, index) {
                          final previsao = previsoes[index];
                          return ForecastCard(previsao: previsao);
                        },
                      ),
                    );
                  }
                  return const Center(
                      child: Text('Nenhuma previsão encontrada.',
                          style: TextStyle(color: Colors.white70)));
                },
              ),
            ),
            const Spacer(), // Adiciona um espaço flexível para empurrar o conteúdo para cima
          ],
        ),
      ),
    );
  }
}

// Widget separado para o card da previsão diária
class ForecastCard extends StatelessWidget {
  final PrevisaoDia previsao;
  const ForecastCard({super.key, required this.previsao});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 70, // Largura um pouco menor
      padding: const EdgeInsets.symmetric(vertical: 8),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.1),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            previsao.diaSemana.substring(0, 3),
            style: const TextStyle(
                color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16),
          ),
          const SizedBox(height: 8),
          BoxedIcon(getWeatherIcon(previsao.condicao),
              color: Colors.white, size: 30),
          const SizedBox(height: 8),
          Text(
            '${previsao.maxima}°/${previsao.minima}°',
            style: const TextStyle(color: Colors.white, fontSize: 16),
          ),
        ],
      ),
    );
  }
}