#include <EEPROM.h>

#define EEPROM_VERSION_ADDR 0
#define EEPROM_BRIGHTNESS_ADDR 1
uint8_t brightness_ = 100;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(57600);
  pinMode(3, OUTPUT);
  pinMode(11, OUTPUT);
  TCCR2A = _BV(COM2A1) | _BV(COM2A0)| _BV(COM2B1) | _BV(WGM20);
  TCCR2B = _BV(CS22);
  brightness_ = EEPROM.read(EEPROM_BRIGHTNESS_ADDR);
  pwm_off();
}


inline void pwm_on(void)
{
  OCR2A = 255 - brightness_;
  OCR2B = brightness_;
}

inline void pwm_off(void)
{
  OCR2A = 255;
  OCR2B = 0;
}

inline void pwm_set(uint8_t val)
{
  brightness_=val;
  pwm_on();  
  EEPROM.write(EEPROM_BRIGHTNESS_ADDR,brightness_);
}

inline void pwm_inc(void)
{
  if (OCR2B < 130)
  {
    OCR2A--;
    OCR2B++;
    brightness_ = OCR2B;
    EEPROM.write(EEPROM_BRIGHTNESS_ADDR,brightness_);
  }
}

inline void pwm_dec(void)
{
  if (OCR2B >0)
  {
    OCR2A++;
    OCR2B--;
    brightness_ = OCR2B;
    EEPROM.write(EEPROM_BRIGHTNESS_ADDR,brightness_);
  }
}


void loop() {
  // put your main code here, to run repeatedly:

  if (Serial.available())
  {
    char cmd = Serial.read();
    switch (cmd) {
      case '0': pwm_off(); Serial.print("off\r\n"); break;
      case '1': pwm_on();  Serial.print("on\r\n"); break;
      case '=': pwm_set(100);  Serial.print("default\r\n"); break;
      case '+': pwm_inc(); Serial.print("pwm = "); Serial.println(OCR2B);  break;
      case '-': pwm_dec(); Serial.print("pwm = "); Serial.println(OCR2B);  break;
      default: Serial.print("error\r\n"); return;
    }
  }
}
