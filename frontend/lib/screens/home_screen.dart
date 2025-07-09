import 'package:flutter/material.dart';
import 'package:weather_icons/weather_icons.dart';
import '../api/api_service.dart';
import '../models/clima_models.dart';
import '../utils/weather_helper.dart';
import 'detail_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  late Future<List<ClimaAtual>> futureClimas;

  @override
  void initState() {
    super.initState();
    futureClimas = fetchClimasAtuais();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Clima das Cidades'),
      ),
      body: FutureBuilder<List<ClimaAtual>>(
        future: futureClimas,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(
                child: CircularProgressIndicator(color: Colors.white));
          } else if (snapshot.hasError) {
            return Center(
                child: Text('Erro: ${snapshot.error}',
                    style: const TextStyle(color: Colors.white70)));
          } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
            return ListView.builder(
              padding: const EdgeInsets.all(8),
              itemCount: snapshot.data!.length,
              itemBuilder: (context, index) {
                final clima = snapshot.data![index];
                final nomeCidadeCorrigido =
                    clima.cidade.toLowerCase() == 'sao paulo'
                        ? 'São Paulo'
                        : clima.cidade;
                return Card(
                  color: Colors.white.withOpacity(0.1),
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(15)),
                  child: ListTile(
                    contentPadding: const EdgeInsets.symmetric(
                        horizontal: 20, vertical: 10),
                    leading: BoxedIcon(getWeatherIcon(clima.condicao),
                        color: Colors.white, size: 28),
                    title: Text(nomeCidadeCorrigido,
                        style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 18,
                            color: Colors.white)),
                    subtitle: Text(clima.condicao,
                        style: TextStyle(color: Colors.white.withOpacity(0.7))),
                    trailing: Text('${clima.temperatura}°C',
                        style: const TextStyle(
                            fontSize: 22,
                            fontWeight: FontWeight.bold,
                            color: Colors.white)),
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (context) => DetailScreen(clima: clima)),
                      );
                    },
                  ),
                );
              },
            );
          }
          return const Center(
              child: Text('Nenhuma cidade encontrada.',
                  style: TextStyle(color: Colors.white70)));
        },
      ),
    );
  }
}