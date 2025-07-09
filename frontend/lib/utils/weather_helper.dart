import 'package:flutter/material.dart';
import 'package:weather_icons/weather_icons.dart';

IconData getWeatherIcon(String condition) {
  condition = condition.toLowerCase();
  if (condition.contains('chuva')) return WeatherIcons.rain;
  if (condition.contains('sol') || condition.contains('limpo')) {
    return WeatherIcons.day_sunny;
  }
  if (condition.contains('nublado')) return WeatherIcons.cloudy;
  if (condition.contains('tempestade')) return WeatherIcons.thunderstorm;
  if (condition.contains('neve')) return WeatherIcons.snow;
  if (condition.contains('névoa') || condition.contains('neblina')) {
    return WeatherIcons.fog;
  }
  return WeatherIcons.na;
}