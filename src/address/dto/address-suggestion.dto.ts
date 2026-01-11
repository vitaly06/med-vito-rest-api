import { ApiProperty } from '@nestjs/swagger';

import { IsNumber, IsOptional, IsString } from 'class-validator';

export class AddressSuggestion {
  @ApiProperty()
  @IsString()
  query: string;

  @ApiProperty({
    required: false,
  })
  @IsOptional()
  @IsNumber()
  limit?: number = 5;
}
