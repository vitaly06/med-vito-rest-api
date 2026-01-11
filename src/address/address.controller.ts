import {
  Body,
  Controller,
  Get,
  HttpException,
  HttpStatus,
  Post,
  Query,
} from '@nestjs/common';

import { AddressService } from './address.service';
import { AddressSuggestion } from './dto/address-suggestion.dto';
import { CreateAddressDto } from './dto/create-address.dto';

@Controller('address')
export class AddressController {
  constructor(private readonly addressService: AddressService) {}

  @Get('suggestions')
  async getSuggestions(@Query() dto: AddressSuggestion) {
    return this.addressService.getSuggestions(dto.query, dto.limit);
  }

  @Post('validate')
  async validateAddress(@Body() dto: CreateAddressDto) {
    const isValid = await this.addressService.validateAddress(dto.address);

    if (!isValid) {
      throw new HttpException(
        'Указанный адрес не найден. Пожалуйста, выберите адрес из предложенных вариантов.',
        HttpStatus.BAD_REQUEST,
      );
    }

    return { valid: true, message: 'Адрес корректен' };
  }
}
