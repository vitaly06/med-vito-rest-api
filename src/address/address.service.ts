import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { ConfigService } from '@nestjs/config';
import { lastValueFrom } from 'rxjs';

@Injectable()
export class AddressService {
  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService,
  ) {}

  async getSuggestions(query: string, limit: number = 5) {
    try {
      const apiKey = this.configService.get<string>('DADATA_API_KEY');

      if (!apiKey) {
        throw new Error('DADATA_API_KEY не найден в конфигурации');
      }

      const response$ = this.httpService.post(
        'https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address',
        {
          query,
          count: limit,
        },
        {
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
            Authorization: `Token ${apiKey}`,
          },
        },
      );

      const response = await lastValueFrom(response$);

      return response.data.suggestions.map((suggestion) => ({
        value: suggestion.value,
      }));
    } catch (error) {
      console.error('DaData API Error:', error.response?.data || error.message);
      throw new HttpException(
        error.response?.data?.message || 'Ошибка получения подсказок адреса',
        error.response?.status || HttpStatus.BAD_GATEWAY,
      );
    }
  }

  async validateAddress(address: string): Promise<boolean> {
    try {
      const suggestions = await this.getSuggestions(address, 1);
      return (
        suggestions.length > 0 &&
        suggestions[0].value.toLowerCase() === address.toLowerCase()
      );
    } catch (error) {
      return false;
    }
  }
}
